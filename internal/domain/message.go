package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageStatus string

const (
	StatusPending MessageStatus = "pending"
	StatusSent    MessageStatus = "sent"
	StatusFailed  MessageStatus = "failed"
)

const MaxMessageLength = 160

var (
	ErrMessageTooLong     = errors.New("message content exceeds maximum length")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
	ErrEmptyContent       = errors.New("message content cannot be empty")
)

type Message struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	Content     string             `json:"content" bson:"content"`
	Status      MessageStatus      `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	SentAt      *time.Time         `json:"sent_at,omitempty" bson:"sent_at,omitempty"`
	MessageID   *string            `json:"message_id,omitempty" bson:"message_id,omitempty"`
}

func (m *Message) Validate() error {
	if m.Content == "" {
		return ErrEmptyContent
	}

	if len(m.Content) > MaxMessageLength {
		return ErrMessageTooLong
	}

	if m.PhoneNumber == "" {
		return ErrInvalidPhoneNumber
	}

	return nil
}

func (m *Message) MarkAsSent(messageID string) {
	now := time.Now()
	m.Status = StatusSent
	m.SentAt = &now
	m.MessageID = &messageID
}

func (m *Message) MarkAsFailed() {
	m.Status = StatusFailed
}
