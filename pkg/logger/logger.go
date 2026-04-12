package logger

import (
	"os"
	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

type Logger struct {
	*logrus.Logger
}

func New() *Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout) // Явно указываем вывод
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return &Logger{l}
}

// Позволяет писать h.logger.WithFields(logger.Fields{...})
func (l *Logger) WithFields(fields Fields) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields(fields))
}

// Позволяет писать h.logger.WithError(err)
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}