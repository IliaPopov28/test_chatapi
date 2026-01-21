FROM golang:1.25-alpine

WORKDIR /app

# Install git and other dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go modules files
COPY go.mod go.sum ./
RUN go mod download

# Install goose
RUN go install github.com/pressly/goose/cmd/goose@v2.1.0+incompatible

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/api

# Expose port
EXPOSE 8080

# Run migrations and start the application
CMD ["sh", "-c", "goose -dir ./migrations postgres \"user=postgres password=password dbname=chatdb host=db port=5432 sslmode=disable\" up && ./main"]
