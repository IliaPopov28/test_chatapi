package model

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Chat struct {
	ID        uint       `json:"id" gorm:"primary_key;auto_increment"`
	Title     string     `json:"title" gorm:"size:200;not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	Messages  []*Message `json:"messages" gorm:"foreignKey:ChatID"`
}

func (c *Chat) BeforeCreate(scope *gorm.Scope) error {
	if c.Title != "" {
		trimmed := strings.TrimSpace(c.Title)
		scope.SetColumn("Title", trimmed)
	}
	return nil
}
