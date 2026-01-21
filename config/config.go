package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит конфигурационные параметры приложения
type Config struct {
	// Порт для запуска HTTP-сервера
	ServerPort string

	// Конфигурация базы данных PostgreSQL
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Таймауты для HTTP-сервера
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Таймаут для graceful shutdown
	ShutdownTimeout time.Duration
}

// LoadConfig загружает конфигурацию из переменных окружения или файла .env
func LoadConfig() (*Config, error) {
	// Попробуем загрузить переменные из файла .env
	godotenv.Load()

	cfg := &Config{
		ServerPort: os.Getenv("SERVER_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	// Значения по умолчанию
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
	if cfg.DBHost == "" {
		cfg.DBHost = "localhost"
	}
	if cfg.DBPort == "" {
		cfg.DBPort = "5432"
	}
	if cfg.DBUser == "" {
		cfg.DBUser = "postgres"
	}
	if cfg.DBPassword == "" {
		cfg.DBPassword = "password"
	}
	if cfg.DBName == "" {
		cfg.DBName = "chatdb"
	}

	// Таймауты по умолчанию
	cfg.ReadTimeout = getDurationEnv("READ_TIMEOUT", 15*time.Second)
	cfg.WriteTimeout = getDurationEnv("WRITE_TIMEOUT", 15*time.Second)
	cfg.IdleTimeout = getDurationEnv("IDLE_TIMEOUT", 60*time.Second)
	cfg.ShutdownTimeout = getDurationEnv("SHUTDOWN_TIMEOUT", 5*time.Second)

	return cfg, nil
}

// getDurationEnv парсит переменную окружения в time.Duration
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
