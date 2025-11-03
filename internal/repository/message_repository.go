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

// SeedSampleData database boşsa örnek mesajlar ekler
func (r *messageRepository) SeedSampleData(ctx context.Context) error {
	// Mesaj sayısını kontrol et
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to count documents: %w", err)
	}

	// Eğer veri varsa seed etme
	if count > 0 {
		return nil
	}

	// Sample mesajlar oluştur
	sampleMessages := []interface{}{
		&domain.Message{
			ID:          primitive.NewObjectID(),
			PhoneNumber: "+905551111111",
			Content:     "Test mesaji 1 - Insider Project",
			Status:      domain.StatusPending,
			CreatedAt:   time.Now(),
		},
		&domain.Message{
			ID:          primitive.NewObjectID(),
			PhoneNumber: "+905552222222",
			Content:     "Test mesaji 2 - Welcome to Insider",
			Status:      domain.StatusPending,
			CreatedAt:   time.Now(),
		},
		&domain.Message{
			ID:          primitive.NewObjectID(),
			PhoneNumber: "+905553333333",
			Content:     "Test mesaji 3 - Siparişiniz hazır",
			Status:      domain.StatusPending,
			CreatedAt:   time.Now(),
		},
	}

	_, err = r.collection.InsertMany(ctx, sampleMessages)
	if err != nil {
		return fmt.Errorf("failed to insert sample data: %w", err)
	}

	return nil
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
