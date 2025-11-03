package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/domain"
	"github.com/UmutcanKalkan/auto-message-dispatcher/pkg/logger"
)

type mockMessageService struct {
	callCount int
}

func (m *mockMessageService) ProcessPendingMessages(ctx context.Context, batchSize int) error {
	m.callCount++
	return nil
}

func (m *mockMessageService) GetSentMessages(ctx context.Context) ([]*domain.Message, error) {
	return nil, nil
}

func (m *mockMessageService) CreateMessage(ctx context.Context, phoneNumber, content string) error {
	return nil
}

func TestScheduler_StartStop(t *testing.T) {
	mock := &mockMessageService{}
	log := logger.New()

	s := NewScheduler(mock, 100*time.Millisecond, 2, log)

	if s.IsRunning() {
		t.Error("should not be running initially")
	}

	s.Start()
	if !s.IsRunning() {
		t.Error("should be running after start")
	}

	time.Sleep(250 * time.Millisecond)

	if mock.callCount < 2 {
		t.Errorf("expected at least 2 calls, got %d", mock.callCount)
	}

	s.Stop()
	if s.IsRunning() {
		t.Error("should not be running after stop")
	}
}

func TestScheduler_RestartCapability(t *testing.T) {
	mock := &mockMessageService{}
	log := logger.New()

	s := NewScheduler(mock, 1*time.Second, 2, log)

	s.Start()
	s.Stop()

	s.Start()
	time.Sleep(100 * time.Millisecond)
	s.Stop()

	if s.IsRunning() {
		t.Error("should be stopped")
	}
}
