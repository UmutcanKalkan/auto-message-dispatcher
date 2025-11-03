package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository interface {
	GetPendingMessages(ctx context.Context, limit int) ([]*domain.Message, error)
	GetSentMessages(ctx context.Context) ([]*domain.Message, error)
	UpdateMessageStatus(ctx context.Context, message *domain.Message) error
	CreateMessage(ctx context.Context, message *domain.Message) error
}

type messageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &messageRepository{
		collection: db.Collection("messages"),
	}
}

func (r *messageRepository) GetPendingMessages(ctx context.Context, limit int) ([]*domain.Message, error) {
	filter := bson.M{"status": domain.StatusPending}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending messages: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var messages []*domain.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}

	return messages, nil
}

func (r *messageRepository) GetSentMessages(ctx context.Context) ([]*domain.Message, error) {
	filter := bson.M{"status": domain.StatusSent}
	opts := options.Find().SetSort(bson.D{{Key: "sent_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query sent messages: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var messages []*domain.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}

	return messages, nil
}

func (r *messageRepository) UpdateMessageStatus(ctx context.Context, message *domain.Message) error {
	filter := bson.M{"_id": message.ID}
	update := bson.M{
		"$set": bson.M{
			"status":     message.Status,
			"sent_at":    message.SentAt,
			"message_id": message.MessageID,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("message not found: %v", message.ID)
	}

	return nil
}

func (r *messageRepository) CreateMessage(ctx context.Context, message *domain.Message) error {
	message.ID = primitive.NewObjectID()
	message.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}
