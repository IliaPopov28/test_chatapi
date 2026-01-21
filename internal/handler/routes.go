package handler

import (
	"net/http"

	"chat-api/internal/service"
)

// NewRouter создает и конфигурирует маршруты для API
func NewRouter(chatService *service.ChatService, messageService *service.MessageService) http.Handler {
	chatHandler := NewChatHandler(chatService, messageService)
	messageHandler := NewMessageHandler(messageService)

	// Создаем мультиплексер для маршрутизации
	mux := http.NewServeMux()

	// Маршруты для чатов
	mux.HandleFunc("/chats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			chatHandler.CreateChat(w, r)
		} else {
			JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/chats/", func(w http.ResponseWriter, r *http.Request) {
		handleChatRoutes(w, r, chatHandler, messageHandler)
	})

	return mux
}

func handleChatRoutes(w http.ResponseWriter, r *http.Request, chatHandler *ChatHandler, messageHandler *MessageHandler) {
	switch {
	case r.URL.Path == "/chats/":
		if r.Method == http.MethodPost {
			chatHandler.CreateChat(w, r)
		} else {
			JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	default:
		// Обработка маршрутов с id
		parts := splitPath(r.URL.Path)

		switch {
		case len(parts) == 2: // /chats/{id}
			switch r.Method {
			case http.MethodGet:
				chatHandler.GetChat(w, r)
			case http.MethodDelete:
				chatHandler.DeleteChat(w, r)
			default:
				JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		case len(parts) == 3 && parts[2] == "messages": // /chats/{id}/messages
			if r.Method == http.MethodPost {
				messageHandler.CreateMessage(w, r)
			} else {
				JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		default:
			JSONResponse(w, http.StatusNotFound, "Endpoint not found")
		}
	}
}

func splitPath(path string) []string {
	var result []string
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if i > start {
				result = append(result, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		result = append(result, path[start:])
	}
	return result
}
