package log

import (
	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

const (
	// for kratos, it need caller skip 2 by logger.Info and filter.
	callerSkipLevelNum = 2
)

func NewLoggerByZap(isTrace bool, logLevelStr string) tmlog.Logger {
	zapLogger := NewZapLogger(viper.GetBool(cli.TraceFlag))

	// warp zap log to logger, it will add caller skip 1
	logger := NewLogger(zapLogger)

	// add caller skip by 2, as warp and level log
	logger = logger.WithCallerSkip(callerSkipLevelNum)

	// process log level for cosmos-sdk, , it will add caller skip 1
	loggerByLevel, err := tmflags.ParseLogLevel(logLevelStr, logger, cfg.DefaultLogLevel())
	if err != nil {
		panic(err)
	}

	return loggerByLevel
}
