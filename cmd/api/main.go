package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"chat-api/config"
	"chat-api/internal/handler"
	"chat-api/internal/repository"
	"chat-api/internal/service"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключение к базе данных
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Инициализация зависимостей
	repo := repository.NewRepository(db)
	chatService := service.NewChatService(repo)
	messageService := service.NewMessageService(repo)

	// Создание маршрутизатора
	mux := handler.NewRouter(chatService, messageService)

	// Создание HTTP-сервера с таймаутами из конфигурации
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Канал для обработки сигналов прерывания
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Запуск сервера в горутине
	go func() {
		log.Printf("Сервер запущен на :%s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидание сигнала прерывания
	<-stopChan

	log.Println("Закрытие сервера...")

	// Создание контекста с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при закрытии сервера: %v", err)
	}

	log.Println("Сервер успешно закрыт")
}
