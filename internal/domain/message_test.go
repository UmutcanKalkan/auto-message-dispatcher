package domain

import (
	"strings"
	"testing"
)

func TestMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		wantErr error
	}{
		{
			name: "valid message",
			message: Message{
				PhoneNumber: "+905551111111",
				Content:     "Test message",
			},
			wantErr: nil,
		},
		{
			name: "empty content",
			message: Message{
				PhoneNumber: "+905551111111",
				Content:     "",
			},
			wantErr: ErrEmptyContent,
		},
		{
			name: "empty phone",
			message: Message{
				PhoneNumber: "",
				Content:     "Test",
			},
			wantErr: ErrInvalidPhoneNumber,
		},
		{
			name: "too long",
			message: Message{
				PhoneNumber: "+905551111111",
				Content:     strings.Repeat("a", 161),
			},
			wantErr: ErrMessageTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessage_MarkAsSent(t *testing.T) {
	msg := &Message{
		PhoneNumber: "+905551111111",
		Content:     "Test",
		Status:      StatusPending,
	}

	messageID := "test-id"
	msg.MarkAsSent(messageID)

	if msg.Status != StatusSent {
		t.Errorf("status = %s, want %s", msg.Status, StatusSent)
	}

	if msg.MessageID == nil || *msg.MessageID != messageID {
		t.Error("messageID not set correctly")
	}

	if msg.SentAt == nil {
		t.Error("sentAt not set")
	}
}
