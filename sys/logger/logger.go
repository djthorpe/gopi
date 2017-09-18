package logger /* import "github.com/djthorpe/gopi/sys/logger" */

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTS

// The level of the logging required
type Level uint

// The configuration for the logger
type Config struct {
	Level  Level
	Path   string
	Append bool
}

// The driver for the logging
type driver struct {
	level  Level
	device *os.File
	mutex  sync.Mutex
}

///////////////////////////////////////////////////////////////////////////////
// CONSTS

const (
	LOG_ANY Level = iota
	LOG_DEBUG2
	LOG_DEBUG
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_FATAL
	LOG_NONE
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register logger
	gopi.RegisterModule(gopi.Module{
		Name:   "sys/logger",
		Type:   gopi.MODULE_TYPE_LOGGER,
		Config: configLogger,
		New:    newLogger,
	})
}

////////////////////////////////////////////////////////////////////////////////
// CONFIG AND NEW

func configLogger(config *gopi.AppConfig) {
	config.AppFlags.FlagString("log", "", "File for logging (default: log to stderr)")
	config.AppFlags.FlagBool("logappend", false, "When writing log to file, append output to end of file")
}

func newLogger(app *gopi.AppInstance) (gopi.Driver, error) {
	path, _ := app.AppFlags.GetString("log")
	append, _ := app.AppFlags.GetBool("logappend")
	return gopi.Open(Config{
		Path:   path,
		Append: append,
		Level:  getLevelForApp(app),
	}, nil)
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open a logger
func (config Config) Open(_ gopi.Logger) (gopi.Driver, error) {
	var err error

	this := new(driver)
	this.level = config.Level

	// Open stderr or a device
	if strings.TrimSpace(config.Path) == "" {
		this.device = os.Stderr
	} else {
		flag := os.O_RDWR | os.O_CREATE
		if config.Append {
			flag |= os.O_APPEND
		}
		if this.device, err = os.OpenFile(config.Path, flag, 0666); err != nil {
			return nil, err
		}
	}
	return this, nil
}

// Close a logger
func (this *driver) Close() error {
	if this.device != nil && this.device != os.Stdout && this.device != os.Stderr {
		return this.device.Close()
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// LOGGING INTERFACE

// Get logging level
func (this *driver) Level() Level {
	return this.level
}

// Set logging level
func (this *driver) SetLevel(level Level) {
	this.level = level
}

func (this *driver) Info(format string, v ...interface{}) {
	if this.level <= LOG_INFO || this.level == LOG_ANY {
		this.log(LOG_INFO, fmt.Sprintf(format, v...))
	}
}

func (this *driver) Debug(format string, v ...interface{}) {
	if this.level <= LOG_DEBUG || this.level == LOG_ANY {
		this.log(LOG_DEBUG, fmt.Sprintf(format, v...))
	}
}

func (this *driver) Debug2(format string, v ...interface{}) {
	if this.level <= LOG_DEBUG2 || this.level == LOG_ANY {
		this.log(LOG_DEBUG2, fmt.Sprintf(format, v...))
	}
}

func (this *driver) Warn(format string, v ...interface{}) {
	if this.level <= LOG_WARN || this.level == LOG_ANY {
		this.log(LOG_WARN, fmt.Sprintf(format, v...))
	}
}

func (this *driver) Error(format string, v ...interface{}) gopi.Error {
	message := fmt.Sprintf(format, v...)
	if this.level <= LOG_ERROR || this.level == LOG_ANY {
		this.log(LOG_ERROR, message)
	}
	return gopi.NewError(errors.New(message))
}

func (this *driver) Fatal(format string, v ...interface{}) gopi.Error {
	message := fmt.Sprintf(format, v...)
	if this.level <= LOG_FATAL || this.level == LOG_ANY {
		this.log(LOG_FATAL, message)
	}
	return gopi.NewError(errors.New(message))
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func getLevelForApp(app *gopi.AppInstance) Level {
	if app.Debug() {
		if app.Verbose() {
			return LOG_DEBUG2
		} else {
			return LOG_DEBUG
		}
	} else if app.Verbose() {
		return LOG_INFO
	}
	return LOG_WARN
}

func (this *driver) log(l Level, message string) {
	if this.device != nil {
		this.mutex.Lock()
		defer this.mutex.Unlock()
		fmt.Fprintf(this.device, "[%v] %v\n", l, message)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (l Level) String() string {
	switch l {
	case LOG_DEBUG2:
		return "DEBUG"
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_INFO:
		return "INFO"
	case LOG_WARN:
		return "WARN"
	case LOG_ERROR:
		return "ERROR"
	case LOG_FATAL:
		return "FATAL"
	default:
		return "[Invalid Level value]"
	}
}

func (this *driver) String() string {
	return fmt.Sprintf("sys.logger{ level=%v }", this.level)
}
