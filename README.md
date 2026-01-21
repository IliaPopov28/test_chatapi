# Chat and Messages API

RESTful API для управления чатами и сообщениями с использованием Go, GORM, PostgreSQL и Docker.

## Технический стек
- **Go 1.25** - язык программирования
- **net/http** - стандартная веб-библиотека Go
- **GORM** - ORM для работы с базой данных
- **PostgreSQL 15** - реляционная база данных
- **Goose** - инструмент для миграций БД
- **Docker** - контейнеризация
- **Testify** - библиотека для тестирования

## Установка и запуск

### Предварительные требования
- Docker и Docker Compose должны быть установлены

### Запуск с помощью Docker Compose
1. Клонируйте репозиторий
2. Перейдите в каталог проекта
3. Запустите сервисы:
   ```bash
   docker-compose up -d
   ```

Приложение будет доступно по адресу http://localhost:8080

### Настройка переменных окружения
Создайте файл `.env` в корне проекта (пример содержимого):
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=chatdb
SERVER_PORT=8080
READ_TIMEOUT=15s
WRITE_TIMEOUT=15s
IDLE_TIMEOUT=60s
SHUTDOWN_TIMEOUT=5s
```

## Методы API

### Чаты
- **Создать чат**: POST /chats
- **Получить чат с сообщениями**: GET /chats/{id}?limit=20
- **Удалить чат**: DELETE /chats/{id}

### Сообщения
- **Отправить сообщение**: POST /chats/{id}/messages

## Валидация
- **title**: 1-200 символов, не пустой, автоматическое обрезание пробелов
- **text**: 1-5000 символов, не пустой
- **limit**: 1-100 (по умолчанию 20)

## Миграции
- Применить миграции:
  ```bash
  docker exec chat-api goose -dir /app/migrations postgres "user=postgres password=password dbname=chatdb host=db port=5432 sslmode=disable" up
  ```
- Откат миграций:
  ```bash
  docker exec chat-api goose -dir /app/migrations postgres "user=postgres password=password dbname=chatdb host=db port=5432 sslmode=disable" down
  ```

## Тестирование
Запустите тесты в Docker:
```bash
docker exec chat-api go test ./tests -v
```

