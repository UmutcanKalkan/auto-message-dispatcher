package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Scheduler SchedulerConfig
	Webhook   WebhookConfig
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	URI    string
	DBName string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type SchedulerConfig struct {
	Interval         time.Duration
	BatchSize        int
	AutoStartEnabled bool
}

type WebhookConfig struct {
	URL        string
	AuthKey    string
	Timeout    time.Duration
	MaxRetries int
	RetryDelay time.Duration
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			URI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
			DBName: getEnv("MONGO_DB", "message_dispatcher"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		Scheduler: SchedulerConfig{
			Interval:         getDurationEnv("SCHEDULER_INTERVAL", 2*time.Minute),
			BatchSize:        getIntEnv("SCHEDULER_BATCH_SIZE", 2),
			AutoStartEnabled: getBoolEnv("SCHEDULER_AUTO_START", true),
		},
		Webhook: WebhookConfig{
			URL:        getEnv("WEBHOOK_URL", ""),
			AuthKey:    getEnv("WEBHOOK_AUTH_KEY", ""),
			Timeout:    getDurationEnv("WEBHOOK_TIMEOUT", 30*time.Second),
			MaxRetries: getIntEnv("WEBHOOK_MAX_RETRIES", 3),
			RetryDelay: getDurationEnv("WEBHOOK_RETRY_DELAY", 1*time.Second),
		},
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.Webhook.URL == "" {
		return fmt.Errorf("WEBHOOK_URL is required")
	}

	if c.Webhook.AuthKey == "" {
		return fmt.Errorf("WEBHOOK_AUTH_KEY is required")
	}

	if c.Scheduler.BatchSize < 1 {
		return fmt.Errorf("SCHEDULER_BATCH_SIZE must be at least 1")
	}

	if c.Database.URI == "" {
		return fmt.Errorf("MONGO_URI is required")
	}

	return nil
}

func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
