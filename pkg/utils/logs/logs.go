package logs

import (
	"fmt"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	logger *zap.SugaredLogger
}

func InitializeLog(logPath string, opt ...LoggerOptions) (Log, error) {
	var loggerOption options
	for i := range opt {
		opt[i](&loggerOption)
	}

	writer := createLogWriter(logPath)

	encoder := getLogEncoder()

	core := zapcore.NewCore(encoder, writer, loggerOption.level)

	logger := zap.New(core).Sugar()

	return Log{
		logger: logger,
	}, nil
}

func createLogWriter(logPath string) zapcore.WriteSyncer {
	fileName := fmt.Sprintf("/%s/logs.log", logPath)

	lumberJackWriter := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     5,
		Compress:   false,
	}

	return zapcore.AddSync(lumberJackWriter)
}

func getLogEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(config)
}

func (l Log) Info(msg string) {
	l.logger.Info(msg)
}

func (l Log) Error(err error) {
	l.logger.Error(err.Error())
}

func (l Log) Fatal(err error) {
	l.logger.Fatal(err.Error())
}