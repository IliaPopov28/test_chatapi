package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"chat-api/config"
	"chat-api/internal/handler"
	"chat-api/internal/model"
	"chat-api/internal/repository"
	"chat-api/internal/service"

	"github.com/stretchr/testify/assert"
)

func setupTestDB() (*repository.Repository, error) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "password",
		DBName:     "chatdb",
	}

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	return repository.NewRepository(db), nil
}

func TestCreateChat(t *testing.T) {
	repo, err := setupTestDB()
	assert.NoError(t, err)

	chatService := service.NewChatService(repo)
	messageService := service.NewMessageService(repo)
	chatHandler := handler.NewChatHandler(chatService, messageService)

	reqBody := `{"title": "Test Chat"}`
	req := httptest.NewRequest("POST", "/chats", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	chatHandler.CreateChat(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Success bool       `json:"success"`
		Data    model.Chat `json:"data"`
		Error   string     `json:"error"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "Test Chat", response.Data.Title)
	assert.Greater(t, response.Data.ID, uint(0))
}

func TestGetChat(t *testing.T) {
	repo, err := setupTestDB()
	assert.NoError(t, err)

	chatService := service.NewChatService(repo)
	messageService := service.NewMessageService(repo)
	chatHandler := handler.NewChatHandler(chatService, messageService)

	// Create a chat first
	testChat := &model.Chat{Title: "Test Chat for Get"}
	err = repo.CreateChat(testChat)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/chats/"+strconv.Itoa(int(testChat.ID)), nil)
	w := httptest.NewRecorder()

	chatHandler.GetChat(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
		Error   string                 `json:"error"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, float64(testChat.ID), response.Data["id"])
	assert.Equal(t, testChat.Title, response.Data["title"])
}

func TestDeleteChat(t *testing.T) {
	repo, err := setupTestDB()
	assert.NoError(t, err)

	chatService := service.NewChatService(repo)
	messageService := service.NewMessageService(repo)
	chatHandler := handler.NewChatHandler(chatService, messageService)

	// Create a chat first
	testChat := &model.Chat{Title: "Test Chat for Delete"}
	err = repo.CreateChat(testChat)
	assert.NoError(t, err)

	req := httptest.NewRequest("DELETE", "/chats/"+strconv.Itoa(int(testChat.ID)), nil)
	w := httptest.NewRecorder()

	chatHandler.DeleteChat(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify chat was deleted
	deletedChat, err := repo.GetChatByID(testChat.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedChat)
}
