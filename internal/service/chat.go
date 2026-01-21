package service

import (
	"chat-api/internal/model"
	"chat-api/internal/repository"
	"chat-api/internal/validator"
	"errors"

	"github.com/jinzhu/gorm"
)

type ChatService struct {
	repo *repository.Repository
}

func NewChatService(repo *repository.Repository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) CreateChat(title string) (*model.Chat, error) {
	if err := validator.ValidateTitle(title); err != nil {
		return nil, err
	}

	chat := &model.Chat{
		Title: title,
	}

	if err := s.repo.CreateChat(chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) GetChatByID(id uint) (*model.Chat, []*model.Message, error) {
	chat, messages, err := s.repo.GetChatWithMessages(id, validator.DefaultLimit)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, errors.New("chat not found")
		}
		return nil, nil, err
	}
	return chat, messages, nil
}

func (s *ChatService) GetChatWithMessages(id uint, limit int) (*model.Chat, []*model.Message, error) {
	validLimit := validator.ValidateLimit(limit)
	chat, messages, err := s.repo.GetChatWithMessages(id, validLimit)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil, errors.New("chat not found")
		}
		return nil, nil, err
	}
	return chat, messages, nil
}

func (s *ChatService) DeleteChat(id uint) error {
	_, _, err := s.GetChatByID(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteChat(id)
}
