package logging

import "github.com/adetunjii/auth-svc/internal/port"

type Logger struct {
	instance SugarLogger
}

var _ port.AppLogger = (*Logger)(nil)

func NewLogger(zapSugarLogger SugarLogger) *Logger {
	return &Logger{
		instance: zapSugarLogger,
	}
}

func (l *Logger) Info(msg string, kv ...any) {
	l.instance.Infow(msg)
}

func (l *Logger) Error(msg string, err error) {
	if err != nil {
		l.instance.Errorw(msg, "error", err.Error())
	} else {
		l.instance.Errorw(msg, "error", "unknown error")
	}
}

func (l *Logger) Fatal(msg string, err error) {
	if err != nil {
		l.instance.Fatalw(msg, "error", err.Error())
	} else {
		l.instance.Errorw(msg, "error", "unknown error")
	}
}
