package handler

import (
	"encoding/json"
	"net/http"

	"github.com/UmutcanKalkan/auto-message-dispatcher/internal/scheduler"
)

type SchedulerHandler struct {
	scheduler *scheduler.Scheduler
}

func NewSchedulerHandler(scheduler *scheduler.Scheduler) *SchedulerHandler {
	return &SchedulerHandler{
		scheduler: scheduler,
	}
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *SchedulerHandler) Start(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.scheduler.IsRunning() {
		h.sendResponse(w, Response{
			Success: true,
			Message: "Scheduler is already running",
		})
		return
	}

	if err := h.scheduler.Start(); err != nil {
		h.sendError(w, "Failed to start scheduler: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendResponse(w, Response{
		Success: true,
		Message: "Scheduler started successfully",
	})
}

func (h *SchedulerHandler) Stop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !h.scheduler.IsRunning() {
		h.sendResponse(w, Response{
			Success: true,
			Message: "Scheduler is not running",
		})
		return
	}

	if err := h.scheduler.Stop(); err != nil {
		h.sendError(w, "Failed to stop scheduler: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendResponse(w, Response{
		Success: true,
		Message: "Scheduler stopped successfully",
	})
}

func (h *SchedulerHandler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	isRunning := h.scheduler.IsRunning()
	status := "stopped"
	if isRunning {
		status = "running"
	}

	h.sendResponse(w, Response{
		Success: true,
		Message: "Scheduler status retrieved successfully",
		Data: map[string]interface{}{
			"status":  status,
			"running": isRunning,
		},
	})
}

func (h *SchedulerHandler) sendResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func (h *SchedulerHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
	})
	if err != nil {
		return
	}
}
