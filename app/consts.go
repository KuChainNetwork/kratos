package app

import (
	"os"
)

const appName = "KratosApp"

var (
	// DefaultCLIHome default home directories for ktscli
	DefaultCLIHome = os.ExpandEnv("$HOME/.ktscli")

	// DefaultNodeHome default home directories for ktsd
	DefaultNodeHome = os.ExpandEnv("$HOME/.ktsd")
)
