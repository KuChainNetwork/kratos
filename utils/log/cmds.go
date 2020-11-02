package log

import (
	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

func NewLoggerByZap(isTrace bool) tmlog.Logger {
	zapLogger := mkZapLogger(viper.GetBool(cli.TraceFlag))

	// process log level for cosmos-sdk
	logLvCfg := viper.GetString("log_level")
	logger, err := tmflags.ParseLogLevel(logLvCfg, NewLogger(zapLogger), cfg.DefaultLogLevel())
	if err != nil {
		panic(err)
	}

	return logger
}
