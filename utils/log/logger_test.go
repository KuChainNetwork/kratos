package log

import (
	"fmt"

	"testing"

	tmlog "github.com/tendermint/tendermint/libs/log"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(NewZapLogger(false)).WithCallerSkip(1) // for warp
	defer logger.Flush()

	var ll tmlog.Logger = logger

	ll.Info("hello", "v", "isOk", "i", 223)
}

func logWarp(l tmlog.Logger, format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...), "warped", "ok")
}

func TestLoggerWarp(t *testing.T) {
	logger := NewLogger(NewZapLogger(false))
	defer logger.Flush()

	var ll tmlog.Logger = logger.WithCallerSkip(2)

	logWarp(ll, "log in there %s", "curr")
}
