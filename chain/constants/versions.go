package constants

import (
	"github.com/tendermint/tendermint/libs/log"
	"gopkg.in/yaml.v2"
)

var (
	KuchainBuildVersion    = ""
	KuchainBuildBranch     = ""
	KuchainBuildTime       = ""
	KuchainBuildSDKVersion = ""
)

// VersionInfo get Version Info
func VersionInfo() []byte {
	ver := struct {
		Version    string `json:"version" yaml:"version"`
		Branch     string `json:"branch" yaml:"branch"`
		BuildTime  string `json:"build_time" yaml:"build_time"`
		SDKVersion string `json:"sdk_version" yaml:"sdk_version"`
	}{
		Version:    KuchainBuildVersion,
		Branch:     KuchainBuildBranch,
		BuildTime:  KuchainBuildTime,
		SDKVersion: KuchainBuildSDKVersion,
	}

	res, _ := yaml.Marshal(ver)
	return res
}

// LogVersion log version info
func LogVersion(logger log.Logger) {
	logger.Info("Kuchain Version",
		"version", KuchainBuildVersion,
		"branch", KuchainBuildBranch,
		"time", KuchainBuildTime,
		"sdkVersion", KuchainBuildSDKVersion)
}
