package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"chat-api/internal/service"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract chat ID from URL - handle both formats:
	// /chats/{id}/messages or /chats/{id}/messages/
	path := strings.Trim(r.URL.Path, "/") // Remove trailing slash
	parts := strings.Split(path, "/")

	var idStr string
	found := false
	for i, part := range parts {
		if part == "chats" && i+1 < len(parts) {
			idStr = parts[i+1]
			found = true
			break
		}
	}

	if !found {
		JSONResponse(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	var request struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		JSONResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	defer r.Body.Close()

	message, err := h.messageService.CreateMessage(uint(id), request.Text)
	if err != nil {
		if err.Error() == "chat not found" {
			JSONResponse(w, http.StatusNotFound, err)
		} else {
			JSONResponse(w, http.StatusBadRequest, err)
		}
		return
	}

	JSONResponse(w, http.StatusCreated, message)
}
