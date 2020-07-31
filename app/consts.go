package app

import (
	"os"
)

const appName = "KuchainApp"

var (
	// DefaultCLIHome default home directories for kucli
	DefaultCLIHome = os.ExpandEnv("$HOME/.kucli")

	// DefaultNodeHome default home directories for kucd
	DefaultNodeHome = os.ExpandEnv("$HOME/.kucd")
)
