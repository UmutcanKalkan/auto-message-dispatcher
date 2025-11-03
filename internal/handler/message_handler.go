package handler

import (
	"encoding/json"
	"net/http"

	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/service"
)

type MessageHandler struct {
	messageService service.MessageService
}

func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

type CreateMessageRequest struct {
	PhoneNumber string `json:"phone_number"`
	Content     string `json:"content"`
}

func (h *MessageHandler) GetSentMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	messages, err := h.messageService.GetSentMessages(r.Context())
	if err != nil {
		h.sendError(w, "Failed to get sent messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendResponse(w, Response{
		Success: true,
		Message: "ok",
		Data: map[string]interface{}{
			"messages": messages,
			"count":    len(messages),
		},
	})
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.messageService.CreateMessage(r.Context(), req.PhoneNumber, req.Content); err != nil {
		h.sendError(w, "Failed to create message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.sendResponse(w, Response{
		Success: true,
		Message: "created",
	})
}

func (h *MessageHandler) sendResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MessageHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
	})
}
