package logger

import (
	"log"
	"os"
)

// Logger интерфейс для логирования
type Logger struct {
	logger *log.Logger
}

// New создает новый экземпляр Logger
func New() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Info записывает информационное сообщение
func (l *Logger) Info(msg string) {
	l.logger.Println("INFO:", msg)
}

// Infof записывает информационное сообщение с форматированием
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Printf("INFO: "+format, v...)
}

// Warn записывает предупреждающее сообщение
func (l *Logger) Warn(msg string) {
	l.logger.Println("WARN:", msg)
}

// Warnf записывает предупреждающее сообщение с форматированием
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Printf("WARN: "+format, v...)
}

// Error записывает сообщение об ошибке
func (l *Logger) Error(msg string) {
	l.logger.Println("ERROR:", msg)
}

// Errorf записывает сообщение об ошибке с форматированием
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Printf("ERROR: "+format, v...)
}
