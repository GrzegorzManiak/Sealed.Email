package helpers

import (
	"log"
)

var Logger *log.Logger

func SetLogger(logger *log.Logger) {
	if logger == nil {
		panic("Logger cannot be nil")
	}

	if Logger != nil {
		panic("Logger already set")
	}

	Logger = logger
}

func GetLogger() *log.Logger {
	if Logger == nil {
		panic("Logger not set")
	}
	return Logger
}

func CreateDebugLogger() *log.Logger {
	return log.New(log.Writer(), "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
