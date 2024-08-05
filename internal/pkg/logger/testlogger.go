package logger

import (
	"go.uber.org/zap"
)

func NewTestLogger() Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
