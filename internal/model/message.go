package model

import (
	"time"
)

type Message struct {
	ID        uint      `json:"id" gorm:"primary_key;auto_increment"`
	ChatID    uint      `json:"chat_id" gorm:"not null"`
	Text      string    `json:"text" gorm:"size:5000;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Chat      *Chat     `json:"-" gorm:"foreignKey:ChatID"`
}
