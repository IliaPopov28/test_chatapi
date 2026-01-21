package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"chat-api/internal/service"
)

type ChatHandler struct {
	chatService    *service.ChatService
	messageService *service.MessageService
}

func NewChatHandler(chatService *service.ChatService, messageService *service.MessageService) *ChatHandler {
	return &ChatHandler{chatService: chatService, messageService: messageService}
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var request struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		JSONResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	defer r.Body.Close()

	chat, err := h.chatService.CreateChat(request.Title)
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, err)
		return
	}

	JSONResponse(w, http.StatusCreated, chat)
}

func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := r.URL.Path[len("/chats/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	// Get limit from query parameters
	limitStr := r.URL.Query().Get("limit")
	var limit int = 20
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			JSONResponse(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}

	chat, messages, err := h.chatService.GetChatWithMessages(uint(id), limit)
	if err != nil {
		if err.Error() == "chat not found" {
			JSONResponse(w, http.StatusNotFound, err)
		} else {
			JSONResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Create response with chat and messages
	response := struct {
		ID        uint                     `json:"id"`
		Title     string                   `json:"title"`
		CreatedAt string                   `json:"created_at"`
		Messages  []map[string]interface{} `json:"messages"`
	}{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// Convert messages to map for response
	for _, msg := range messages {
		response.Messages = append(response.Messages, map[string]interface{}{
			"id":         msg.ID,
			"chat_id":    msg.ChatID,
			"text":       msg.Text,
			"created_at": msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	JSONResponse(w, http.StatusOK, response)
}

func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := r.URL.Path[len("/chats/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		JSONResponse(w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	err = h.chatService.DeleteChat(uint(id))
	if err != nil {
		if err.Error() == "chat not found" {
			JSONResponse(w, http.StatusNotFound, err)
		} else {
			JSONResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
