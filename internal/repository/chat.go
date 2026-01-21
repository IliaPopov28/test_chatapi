package repository

import (
	"chat-api/internal/model"
)

func (r *Repository) CreateChat(chat *model.Chat) error {
	return r.db.Create(chat).Error
}

func (r *Repository) ChatExists(id uint) bool {
	var count int
	r.db.Model(&model.Chat{}).Where("id = ?", id).Count(&count)
	return count > 0
}

func (r *Repository) GetChatByID(id uint) (*model.Chat, error) {
	var chat model.Chat
	if err := r.db.First(&chat, id).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *Repository) DeleteChat(id uint) error {
	if err := r.db.Where("id = ?", id).Delete(&model.Chat{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetChatWithMessages(id uint, limit int) (*model.Chat, []*model.Message, error) {
	var chat model.Chat
	if err := r.db.First(&chat, id).Error; err != nil {
		return nil, nil, err
	}

	var messages []*model.Message
	if err := r.db.Where("chat_id = ?", id).
		Order("created_at desc").
		Limit(limit).
		Find(&messages).Error; err != nil {
		return nil, nil, err
	}

	return &chat, messages, nil
}
