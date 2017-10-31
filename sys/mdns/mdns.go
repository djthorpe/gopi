/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mdsn /* import "github.com/djthorpe/gopi/sys/mdns" */

import (
	"fmt"

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
func (config Config) Open(log gopi.Logger) (gopi.Driver, error) {

	this := new(driver)
	this.log = log

	// TODO

	// success
	return this, nil
}

// Close a logger
func (this *driver) Close() error {
	// TODO
	return nil
}

func (this *driver) String() string {
	return fmt.Sprintf("sys.mdns{ }")
}
