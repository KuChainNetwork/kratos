package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger create a logger by zap
func NewZapLogger(isDebug bool) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "logTm",
		LevelKey:       "logLv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "logMsg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.EpochNanosTimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if isDebug {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	config := zap.NewProductionConfig()

	if isDebug {
		config = zap.NewDevelopmentConfig()
	}

	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel) // most small
	config.EncoderConfig = encoderConfig
	config.Development = isDebug

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("zap logger build err by %s", err.Error()))
	}

	return logger
}
