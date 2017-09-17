package mock /* import "github.com/djthorpe/gopi/sys/default/logger" */

import (
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register logger
	gopi.RegisterModule(gopi.Module{
		Name:   "default/logger",
		Type:   gopi.MODULE_TYPE_LOGGER,
		Config: configLogger,
		New:    newLogger,
	})
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func configLogger(config *gopi.AppConfig) {
	config.AppFlags.FlagString("log", "", "File for logging (default: log to stderr)")
}

func newLogger(app *gopi.AppInstance) (gopi.Driver, error) {
	logger, err := getLoggerForApp(app)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(getLevelForApp(app))
	return logger, nil
}

func getLoggerForApp(app *gopi.AppInstance) (*util.LoggerDevice, error) {
	file, exists := app.AppFlags.GetString("log")
	if exists {
		return util.Logger(util.FileLogger{Filename: file, Append: false})
	} else {
		return util.Logger(util.StderrLogger{})
	}
}

func getLevelForApp(app *gopi.AppInstance) util.LogLevel {
	if app.Debug() {
		if app.Verbose() {
			return util.LOG_DEBUG2
		} else {
			return util.LOG_DEBUG
		}
	}
	if app.Verbose() {
		return util.LOG_INFO
	} else {
		return util.LOG_WARN
	}
}
