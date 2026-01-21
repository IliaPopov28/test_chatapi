package service

import (
	"chat-api/internal/model"
	"chat-api/internal/repository"
	"chat-api/internal/validator"
	"errors"
)

type MessageService struct {
	repo *repository.Repository
}

func NewMessageService(repo *repository.Repository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) CreateMessage(chatID uint, text string) (*model.Message, error) {
	if err := validator.ValidateText(text); err != nil {
		return nil, err
	}

	// Check if chat exists
	if !s.repo.ChatExists(chatID) {
		return nil, errors.New("chat not found")
	}

	message := &model.Message{
		ChatID: chatID,
		Text:   text,
	}

	if err := s.repo.CreateMessage(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageService) GetMessagesByChatID(chatID uint, limit int) ([]*model.Message, error) {
	validLimit := validator.ValidateLimit(limit)
	return s.repo.GetMessagesByChatID(chatID, validLimit)
}
