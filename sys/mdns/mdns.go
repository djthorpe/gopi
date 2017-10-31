/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mdsn /* import "github.com/djthorpe/gopi/sys/mdns" */

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/third_party/zeroconf"
)

////////////////////////////////////////////////////////////////////////////////
// STRUCTS

// The configuration
type Config struct {
	Domain string
}

// The driver for the logging
type driver struct {
	log     gopi.Logger
	servers []*zeroconf.Server
}

///////////////////////////////////////////////////////////////////////////////
// CONSTS

const (
	MDNS_DOMAIN = "local."
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register logger
	gopi.RegisterModule(gopi.Module{
		Name:   "sys/mdns",
		Type:   gopi.MODULE_TYPE_MDNS,
		Config: configDriver,
		New:    newDriver,
	})
}

////////////////////////////////////////////////////////////////////////////////
// MODULE CONFIG AND NEW

func configDriver(config *gopi.AppConfig) {
	config.AppFlags.FlagString("mdns-domain", MDNS_DOMAIN, "mDNS Network Domain")
}

func newDriver(app *gopi.AppInstance) (gopi.Driver, error) {
	domain, _ := app.AppFlags.GetString("mdns-domain")
	return gopi.Open(Config{
		Domain: domain,
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
