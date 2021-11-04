package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func New(isDebug bool) *Logger {

	logger := logrus.New()

	if isDebug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return &Logger{logger: logger}
}
