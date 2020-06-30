package log

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	l := mkZapLogger(false)
	logger := NewLogger(l)
	defer logger.Flush()

	logger.Info("hello", "v", "isOk", "i", 223)
}
