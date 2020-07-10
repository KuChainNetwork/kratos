package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	chainCfg "github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
)

func PersistentPreRunEFn(context *server.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == version.Cmd.Name() {
			return nil
		}

		zapLogger := mkZapLogger(viper.GetBool(cli.TraceFlag))

		// process log level for cosmos-sdk
		logLvCfg := viper.GetString("log_level")
		logger, err := tmflags.ParseLogLevel(logLvCfg, NewLogger(zapLogger), cfg.DefaultLogLevel())
		if err != nil {
			return err
		}

		context.Logger = logger.With("module", "main")
		context.Config = chainCfg.DefaultConfig()
		return nil
	}
}

func mkZapLogger(isDebug bool) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "tm",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
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

	config := zap.NewDevelopmentConfig()

	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel) // most small
	config.EncoderConfig = encoderConfig
	config.Development = isDebug

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("zap logger build err by %s", err.Error()))
	}

	return logger.WithOptions(zap.AddCallerSkip(2))
}
