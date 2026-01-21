package repository

import (
	"fmt"

	"chat-api/config"
	"chat-api/internal/model"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Auto-migrate models (for development)
	db.AutoMigrate(&model.Chat{}, &model.Message{})

	// Enable logging in development
	db.LogMode(true)

	return db, nil
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Close() error {
	return r.db.Close()
}
