package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/config"
	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/handler"
	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/middleware"
	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/repository"
	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/scheduler"
	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/service"
	"github.com/UmutcanKalkan/auto-message-dispatcher/pkg/database"
	"github.com/UmutcanKalkan/auto-message-dispatcher/pkg/logger"
	"github.com/UmutcanKalkan/auto-message-dispatcher/pkg/redis"
)

func main() {
	log := logger.New()
	log.Info("Starting application...")

	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	log.Info("Connecting to MongoDB...")
	db, err := database.NewMongoDB(cfg.Database.URI, cfg.Database.DBName)
	if err != nil {
		log.Error("Failed to connect to MongoDB: %v", err)
		os.Exit(1)
	}
	log.Info("MongoDB connected successfully")

	log.Info("Connecting to Redis...")
	redisClient, err := redis.NewRedisClient(
		cfg.Redis.GetRedisAddr(),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Error("Failed to connect to Redis: %v", err)
		os.Exit(1)
	}
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {

		}
	}(redisClient)
	log.Info("Redis connected successfully")

	messageRepo := repository.NewMessageRepository(db)

	webhookClient := service.NewWebhookClient(
		cfg.Webhook.URL,
		cfg.Webhook.AuthKey,
		cfg.Webhook.Timeout,
	)

	messageService := service.NewMessageService(
		messageRepo,
		webhookClient,
		redisClient,
		log,
	)

	schedule := scheduler.NewScheduler(
		messageService,
		cfg.Scheduler.Interval,
		cfg.Scheduler.BatchSize,
		log,
	)

	if cfg.Scheduler.AutoStartEnabled {
		log.Info("Auto-starting scheduler...")
		if err := schedule.Start(); err != nil {
			log.Error("Failed to start scheduler: %v", err)
		}
	}

	schedulerHandler := handler.NewSchedulerHandler(schedule)
	messageHandler := handler.NewMessageHandler(messageService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/scheduler/start", schedulerHandler.Start)
	mux.HandleFunc("/api/scheduler/stop", schedulerHandler.Stop)
	mux.HandleFunc("/api/scheduler/status", schedulerHandler.Status)

	mux.HandleFunc("/api/messages/sent", messageHandler.GetSentMessages)
	mux.HandleFunc("/api/messages", messageHandler.CreateMessage)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./api/swagger.yaml")
	})

	cors := middleware.CORS(mux)
	cors = middleware.Logging(log)(cors)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      cors,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Info("Starting HTTP server on port: %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	log.Info("Server started successfully on port: %s", cfg.Server.Port)
	log.Info("API endpoints:")
	log.Info("  POST   /api/scheduler/start")
	log.Info("  POST   /api/scheduler/stop")
	log.Info("  GET    /api/scheduler/status")
	log.Info("  GET    /api/messages/sent")
	log.Info("  POST   /api/messages")
	log.Info("  GET    /health")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	if schedule.IsRunning() {
		log.Info("Stopping scheduler...")
		if err := schedule.Stop(); err != nil {
			log.Error("Failed to stop scheduler: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Info("Server stopped gracefully")
}
