package mock /* import "github.com/djthorpe/gopi/sys/default/logger" */

import (
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register logger
	registerLoggerFlags(gopi.RegisterModule(gopi.Module{
		Name: "default/logger",
		Type: gopi.MODULE_TYPE_LOGGER,
		New:  newLogger,
	}))
}

func registerLoggerFlags(flags *gopi.Flags) {
	flags.FlagString("log", "", "File for logging (default: log to stderr)")
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func newLogger(config *gopi.AppConfig, _ gopi.Logger) (gopi.Driver, error) {
	logger, err := getLoggerForConfig(config)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(getLevelForConfig(config))
	return logger, nil
}

func getLoggerForConfig(config *gopi.AppConfig) (*util.LoggerDevice, error) {
	file, exists := config.AppFlags.GetString("log")
	if exists {
		return util.Logger(util.FileLogger{Filename: file, Append: false})
	} else {
		return util.Logger(util.StderrLogger{})
	}
}

func getLevelForConfig(config *gopi.AppConfig) util.LogLevel {
	if config.Debug {
		if config.Verbose {
			return util.LOG_DEBUG2
		} else {
			return util.LOG_DEBUG
		}
	}
	if config.Verbose {
		return util.LOG_INFO
	} else {
		return util.LOG_WARN
	}
}
