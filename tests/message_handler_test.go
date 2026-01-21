package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"chat-api/internal/handler"
	"chat-api/internal/model"
	"chat-api/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	repo, err := setupTestDB()
	assert.NoError(t, err)

	messageService := service.NewMessageService(repo)
	messageHandler := handler.NewMessageHandler(messageService)

	testChat := &model.Chat{Title: "Test Chat for Messages"}
	err = repo.CreateChat(testChat)
	assert.NoError(t, err)

	reqBody := `{"text": "Hello, world!"}`
	req := httptest.NewRequest("POST", "/chats/"+strconv.Itoa(int(testChat.ID))+"/messages", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	messageHandler.CreateMessage(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Success bool          `json:"success"`
		Data    model.Message `json:"data"`
		Error   string        `json:"error"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "Hello, world!", response.Data.Text)
	assert.Equal(t, testChat.ID, response.Data.ChatID)
	assert.Greater(t, response.Data.ID, uint(0))
	assert.NotZero(t, response.Data.CreatedAt)
}

func TestCreateMessageInNonExistentChat(t *testing.T) {
	repo, err := setupTestDB()
	assert.NoError(t, err)

	messageService := service.NewMessageService(repo)
	messageHandler := handler.NewMessageHandler(messageService)

	reqBody := `{"text": "This chat does not exist!"}`
	req := httptest.NewRequest("POST", "/chats/999/messages", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	messageHandler.CreateMessage(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
		Error   string      `json:"error"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response.Success)
	assert.Equal(t, "chat not found", response.Error)
}

func TestCreateMessageWithInvalidText(t *testing.T) {
	repo, err := setupTestDB()
	assert.NoError(t, err)

	messageService := service.NewMessageService(repo)
	messageHandler := handler.NewMessageHandler(messageService)

	testChat := &model.Chat{Title: "Test Chat for Invalid Messages"}
	err = repo.CreateChat(testChat)
	assert.NoError(t, err)

	reqBody := `{"text": ""}`
	req := httptest.NewRequest("POST", "/chats/"+strconv.Itoa(int(testChat.ID))+"/messages", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	messageHandler.CreateMessage(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
		Error   string      `json:"error"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "text cannot be empty")
}
