/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package logger

import (
	"errors"
	"fmt"
	"log/syslog"
	"os"
	"strings"
	"sync"
	"time"

	// Frameworks
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
	Syslog string
	Tag    string
}

// The driver for the logging
type driver struct {
	level  Level
	device *os.File
	syslog *syslog.Writer
	mutex  sync.Mutex
	delta  time.Time
	tag    string
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

const (
	DELTA_TIMESTAMP_SECS = 60
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
	config.AppFlags.FlagString("log.file", "", "Log to syslog facility (user,daemon,local0...local7) or file (default: log to stderr)")
	config.AppFlags.FlagString("log.tag", "", "Tag for logging (default: name of application)")
	config.AppFlags.FlagBool("log.append", false, "When writing log to file, append output to end of file")
}

func newLogger(app *gopi.AppInstance) (gopi.Driver, error) {
	path, _ := app.AppFlags.GetString("log.file")
	append, _ := app.AppFlags.GetBool("log.append")
	tag, exists := app.AppFlags.GetString("log.tag")
	if exists == false {
		tag = app.AppFlags.Name()
	}
	return gopi.Open(Config{
		Path:   path,
		Append: append,
		Level:  getLevelForApp(app),
		Tag:    tag,
	}, nil)
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open a logger
func (config Config) Open(_ gopi.Logger) (gopi.Driver, error) {
	this := new(driver)
	this.level = config.Level
	this.tag = config.Tag

	if facility, err := getSyslogPriority(config.Path); err != nil && err != gopi.ErrBadParameter {
		// Unknown syslog error
		return nil, err
	} else if err == nil {
		// Syslog facility
		if syslog, err := syslog.New(facility, config.Tag); err != nil {
			return nil, err
		} else {
			this.syslog = syslog
		}
	} else if strings.TrimSpace(config.Path) == "" {
		// Stderr logging
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
	if this.syslog != nil {
		if err := this.syslog.Close(); err != nil {
			return err
		}
	}
	if this.device != nil && this.device != os.Stdout && this.device != os.Stderr {
		if err := this.device.Close(); err != nil {
			return err
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// LOGGING INTERFACE

// Level gets logging level
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

func (this *driver) Error(format string, v ...interface{}) error {
	message := fmt.Sprintf(format, v...)
	if this.level <= LOG_ERROR || this.level == LOG_ANY {
		this.log(LOG_ERROR, message)
	}
	return errors.New(message)
}

func (this *driver) Fatal(format string, v ...interface{}) error {
	message := fmt.Sprintf(format, v...)
	if this.level <= LOG_FATAL || this.level == LOG_ANY {
		this.log(LOG_FATAL, message)
	}
	return errors.New(message)
}

func (this *driver) IsDebug() bool {
	return (this.level == LOG_DEBUG || this.level == LOG_DEBUG2)
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

func getSyslogPriority(value string) (syslog.Priority, error) {
	switch value {
	case "user":
		return syslog.LOG_USER, nil
	case "daemon":
		return syslog.LOG_DAEMON, nil
	case "local0":
		return syslog.LOG_LOCAL0, nil
	case "local1":
		return syslog.LOG_LOCAL1, nil
	case "local2":
		return syslog.LOG_LOCAL2, nil
	case "local3":
		return syslog.LOG_LOCAL3, nil
	case "local4":
		return syslog.LOG_LOCAL4, nil
	case "local5":
		return syslog.LOG_LOCAL5, nil
	case "local6":
		return syslog.LOG_LOCAL6, nil
	case "local7":
		return syslog.LOG_LOCAL7, nil
	default:
		return 0, gopi.ErrBadParameter
	}
}

func (this *driver) log(l Level, message string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.device != nil {
		now := time.Now()
		if this.delta.IsZero() || this.delta.Add(time.Second*DELTA_TIMESTAMP_SECS).Before(now) {
			this.delta = now
			fmt.Fprintf(this.device, "== %v %v ==\n", this.delta.Format(time.RFC3339), this.tag)
		}
		fmt.Fprintf(this.device, "[%v] %v\n", l, message)
	}
	if this.syslog != nil {
		switch l {
		case LOG_DEBUG2, LOG_DEBUG:
			this.syslog.Debug(message)
		case LOG_WARN:
			this.syslog.Warning(message)
		case LOG_INFO:
			this.syslog.Info(message)
		case LOG_ERROR:
			this.syslog.Err(message)
		case LOG_FATAL:
			this.syslog.Crit(message)
		}
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
