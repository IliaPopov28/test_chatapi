package repository

import (
	"chat-api/internal/model"
)

func (r *Repository) CreateMessage(message *model.Message) error {
	return r.db.Create(message).Error
}

func (r *Repository) GetMessagesByChatID(chatID uint, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	if err := r.db.Where("chat_id = ?", chatID).
		Order("created_at desc").
		Limit(limit).
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
