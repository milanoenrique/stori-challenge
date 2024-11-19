package logs

import "go.uber.org/zap/zapcore"

type options struct {
	level zapcore.Level
}

type LoggerOptions func(*options)

type LoggerLevel string

const (
	InfoLevel  LoggerLevel = "info"
	DebugLevel LoggerLevel = "debug"
	ErrorLevel LoggerLevel = "error"
)

func SetLoggerLevel(level LoggerLevel) LoggerOptions {
	l := zapcore.InfoLevel

	switch level {
	case "info":
		l = zapcore.DebugLevel
	case "error":
		l = zapcore.ErrorLevel
	}

	return func(o *options) {
		o.level = l
	}
}