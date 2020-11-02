package log

import (
	"fmt"

	tmlog "github.com/tendermint/tendermint/libs/log"
	"go.uber.org/zap"
)

var _ tmlog.Logger = &Logger{}

// Logger imp logger for tendermint
type Logger struct {
	logger *zap.Logger
}

// NewLogger create a logger by zap.Logger
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func genFields4Log(keyvals ...interface{}) []zap.Field {
	if len(keyvals) == 0 {
		return []zap.Field{}
	}

	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "MissingValue")
	}

	res := make([]zap.Field, 0, len(keyvals)/2)
	for i := 0; i < (len(keyvals) / 2); i++ {
		res = append(res, keyval2Field(keyvals[i*2], keyvals[i*2+1]))
	}

	return res
}

func keyval2Field(key, value interface{}) zap.Field {
	return zap.Any(interface2String(key), value)
}

func (l Logger) Flush() error {
	return l.logger.Sync()
}

// Debug imp for tmlog.Logger
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.logger.Debug(msg, genFields4Log(keyvals...)...)
}

// Info imp for tmlog.Logger
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.logger.Info(msg, genFields4Log(keyvals...)...)
}

// Error imp for tmlog.Logger
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.logger.Error(msg, genFields4Log(keyvals...)...)
}

func interface2String(key interface{}) string {
	str, _ := isValueCanString(key)
	return str
}

func isValueCanString(v interface{}) (string, bool) {
	keyStr, ok := v.(string)
	if ok {
		return keyStr, true
	}

	keyCan2String, ok := v.(fmt.Stringer)
	if ok {
		return keyCan2String.String(), true
	}

	return fmt.Sprintf("%v", v), false
}

// With imp for tmlog.Logger
func (l *Logger) With(keyvals ...interface{}) tmlog.Logger {
	if len(keyvals) <= 2 || len(keyvals)%2 != 0 {
		return l
	}

	fields := genFields4Log(keyvals...)

	return &Logger{
		logger: l.logger.With(fields...),
	}
}

// WithCallerSkip if warp the log, to add the caller skip to show the right caller
func (l *Logger) WithCallerSkip(skip int) *Logger {
	l.logger = l.logger.WithOptions(zap.AddCallerSkip(skip))
	return l
}
