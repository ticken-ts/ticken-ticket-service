package config

import (
	"log"
	"os"
	"sync"
)

type TickenLogger interface {
	Info(message string)
	Debug(message string)
	Error(message string)
}

type tickenLogger struct {
	debug *log.Logger
	info  *log.Logger
	error *log.Logger
}

var logger *tickenLogger
var once sync.Once

func GetTickenLogger() TickenLogger {
	once.Do(func() {
		logger = createLogger()
	})
	return logger
}

func (logger *tickenLogger) Info(message string) {
	logger.info.Println(message)
}

func (logger *tickenLogger) Debug(message string) {
	logger.debug.Println(message)
}

func (logger *tickenLogger) Error(message string) {
	logger.error.Println(message)
}

func createLogger() *tickenLogger {
	return &tickenLogger{
		info:  log.New(os.Stdout, "INFO", log.Ldate|log.Ltime|log.Lshortfile),
		debug: log.New(os.Stdout, "DEBUG", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(os.Stdout, "DEBUG", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
