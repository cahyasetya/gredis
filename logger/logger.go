package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

func InitLogger() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return err
	}

	Log = logger
	Sugar = logger.Sugar()
	zap.ReplaceGlobals(logger)

	return nil
}

// Sync flushes any buffered log entries
func Sync() error {
	return Log.Sync()
}
