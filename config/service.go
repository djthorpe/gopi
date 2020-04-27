/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package config

import (
	"os"
	"path/filepath"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *flagset) setServiceFlags() {
	if hostname, err := os.Hostname(); err == nil {
		if strings.HasSuffix(hostname, ".") == false {
			// Make it fully qualified
			hostname += "."
		}
		this.SetString("host", gopi.FLAG_NS_SERVICE, hostname)
	}
	if executable, err := os.Executable(); err == nil {
		this.SetString("name", gopi.FLAG_NS_SERVICE, filepath.Base(executable))
	}
	// Set service type as gopi
	this.SetString("service", gopi.FLAG_NS_SERVICE, "gopi")
}
