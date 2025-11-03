package service

import (
	"context"
	"fmt"

	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/domain"
	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/repository"
	"github.com/UmutcanKalkan/auto-message-dispatcher/pkg/logger"
	"github.com/UmutcanKalkan/auto-message-dispatcher/pkg/redis"
)

type MessageService interface {
	ProcessPendingMessages(ctx context.Context, batchSize int) error
	GetSentMessages(ctx context.Context) ([]*domain.Message, error)
	CreateMessage(ctx context.Context, phoneNumber, content string) error
}

type messageService struct {
	repo          repository.MessageRepository
	webhookClient WebhookClient
	redisClient   *redis.Client
	logger        *logger.Logger
}

func NewMessageService(
	repo repository.MessageRepository,
	webhookClient WebhookClient,
	redisClient *redis.Client,
	logger *logger.Logger,
) MessageService {
	return &messageService{
		repo:          repo,
		webhookClient: webhookClient,
		redisClient:   redisClient,
		logger:        logger,
	}
}

func (s *messageService) ProcessPendingMessages(ctx context.Context, batchSize int) error {
	messages, err := s.repo.GetPendingMessages(ctx, batchSize)
	if err != nil {
		return fmt.Errorf("failed to get pending messages: %w", err)
	}

	if len(messages) == 0 {
		s.logger.Info("No pending messages to process")
		return nil
	}

	s.logger.Info("Processing %d pending messages", len(messages))

	for _, msg := range messages {
		if err := s.sendMessage(ctx, msg); err != nil {
			s.logger.Error("Failed to send message ID %d: %v", msg.ID, err)
			msg.MarkAsFailed()
			if updateErr := s.repo.UpdateMessageStatus(ctx, msg); updateErr != nil {
				s.logger.Error("Failed to update message status: %v", updateErr)
			}
			continue
		}

		s.logger.Info("Successfully sent message ID %d to %s", msg.ID, msg.PhoneNumber)
	}

	return nil
}

func (s *messageService) sendMessage(ctx context.Context, msg *domain.Message) error {
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("message validation failed: %w", err)
	}

	resp, err := s.webhookClient.SendMessage(ctx, msg.PhoneNumber, msg.Content)
	if err != nil {
		return fmt.Errorf("webhook request failed: %w", err)
	}

	msg.MarkAsSent(resp.MessageID)

	if err := s.repo.UpdateMessageStatus(ctx, msg); err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	if s.redisClient != nil && msg.SentAt != nil {
		if err := s.redisClient.CacheSentMessage(ctx, resp.MessageID, *msg.SentAt); err != nil {
			s.logger.Error("Failed to cache message to Redis: %v", err)
		} else {
			s.logger.Info("Cached message %s to Redis", resp.MessageID)
		}
	}

	return nil
}

func (s *messageService) GetSentMessages(ctx context.Context) ([]*domain.Message, error) {
	messages, err := s.repo.GetSentMessages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent messages: %w", err)
	}
	return messages, nil
}

func (s *messageService) CreateMessage(ctx context.Context, phoneNumber, content string) error {
	message := &domain.Message{
		PhoneNumber: phoneNumber,
		Content:     content,
		Status:      domain.StatusPending,
	}

	if err := message.Validate(); err != nil {
		return fmt.Errorf("message validation failed: %w", err)
	}

	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}
