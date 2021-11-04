package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger2 struct {
	logger *logrus.Logger
}

func New2(isDebug bool) *Logger2 {
	logLevel := logrus.InfoLevel
	if isDebug {
		logLevel = logrus.DebugLevel
	}

	var logger = logrus.New()
	logger.SetLevel(logLevel)

	return &Logger2{logger: logger}
}
