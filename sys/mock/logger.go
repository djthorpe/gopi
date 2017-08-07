package mock /* import "github.com/djthorpe/gopi/sys/mock" */

import (
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register logger
	registerLoggerFlags(gopi.RegisterModule(gopi.Module{
		Name: "mock/logger",
		Type: gopi.MODULE_TYPE_LOGGER,
		New:  newLogger,
	}))
}

func registerLoggerFlags(flags *util.Flags) {
	flags.FlagString("log", "", "File for logging (default: log to stderr)")
	flags.FlagBool("verbose", false, "Log verbosely")
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
	file, exists := config.Flags.GetString("log")
	if exists {
		return util.Logger(util.FileLogger{Filename: file, Append: false})
	} else {
		return util.Logger(util.StderrLogger{})
	}
}

func getLevelForConfig(config *gopi.AppConfig) util.LogLevel {
	verbose, _ := config.Flags.GetBool("verbose")
	if config.Debug {
		if verbose {
			return util.LOG_DEBUG2
		} else {
			return util.LOG_DEBUG
		}
	}
	if verbose {
		return util.LOG_INFO
	} else {
		return util.LOG_WARN
	}
}
