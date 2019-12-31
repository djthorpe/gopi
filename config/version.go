/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package config

import (
	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	GitTag      string
	GitBranch   string
	GitHash     string
	GoBuildTime string
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *flagset) setVersionFlags() {
	flag := false
	if GitTag != "" {
		this.SetString("tag", gopi.FLAG_NS_VERSION, GitTag)
		flag = true
	}
	if GitBranch != "" {
		this.SetString("branch", gopi.FLAG_NS_VERSION, GitBranch)
		flag = true
	}
	if GitHash != "" {
		this.SetString("hash", gopi.FLAG_NS_VERSION, GitHash)
		flag = true
	}
	if GoBuildTime != "" {
		this.SetString("buildtime", gopi.FLAG_NS_VERSION, GoBuildTime)
		flag = true
	}
	if flag {
		this.FlagBool("version", false, "Display version information")
	}
}
