package configuration

import (
	"path"
	"runtime"

	"github.com/JesusIslam/salestock/logger"
)

var DefaultDirectory string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		logger.Fatal("Failed to get runtime caller")
	}

	DefaultDirectory = path.Dir(filename) + "/../default_configurations"
}
