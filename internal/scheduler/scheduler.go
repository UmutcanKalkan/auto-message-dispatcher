package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/service"
	"github.com/UmutcanKalkan/auto-message-dispatcher/pkg/logger"
)

// Scheduler manages message sending at specified intervals
type Scheduler struct {
	messageService service.MessageService
	interval       time.Duration
	batchSize      int
	logger         *logger.Logger

	mu        sync.Mutex
	running   bool
	stopChan  chan struct{}
	doneChan  chan struct{}
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func NewScheduler(
	messageService service.MessageService,
	interval time.Duration,
	batchSize int,
	logger *logger.Logger,
) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		messageService: messageService,
		interval:       interval,
		batchSize:      batchSize,
		logger:         logger,
		stopChan:       make(chan struct{}),
		doneChan:       make(chan struct{}),
		ctx:            ctx,
		cancelCtx:      cancel,
	}
}

func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		s.logger.Info("Scheduler is already running")
		return nil
	}

	s.stopChan = make(chan struct{})
	s.doneChan = make(chan struct{})
	s.ctx, s.cancelCtx = context.WithCancel(context.Background())

	s.running = true
	s.logger.Info("Starting scheduler with interval: %v, batch size: %d", s.interval, s.batchSize)

	go s.run()

	return nil
}

func (s *Scheduler) Stop() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		s.logger.Info("Scheduler is not running")
		return nil
	}
	s.mu.Unlock()

	s.logger.Info("Stopping scheduler...")

	close(s.stopChan)
	s.cancelCtx()
	<-s.doneChan

	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	s.logger.Info("Scheduler stopped")
	return nil
}

func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *Scheduler) run() {
	defer close(s.doneChan)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Info("Scheduler loop started")

	// Process first batch immediately
	s.processBatch()

	for {
		select {
		case <-ticker.C:
			s.processBatch()
		case <-s.stopChan:
			s.logger.Info("Scheduler received stop signal")
			return
		}
	}
}

func (s *Scheduler) processBatch() {
	s.logger.Info("Processing batch of %d messages", s.batchSize)

	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Minute)
	defer cancel()

	if err := s.messageService.ProcessPendingMessages(ctx, s.batchSize); err != nil {
		s.logger.Error("Failed to process pending messages: %v", err)
	}
}
