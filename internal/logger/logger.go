package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	Log  *zap.Logger
	once sync.Once
)

// InitLogger initializes structured logging
func InitLogger() {
	once.Do(func() { // Ensures only one instance is created
		var err error
		Log, err = zap.NewProduction()
		if err != nil {
			panic("failed to initialize logger")
		}
	})
}

// GetLogger ensures the logger is initialized before use
func GetLogger() *zap.Logger {
	if Log == nil {
		InitLogger() // Fallback if not initialized
	}
	return Log
}
