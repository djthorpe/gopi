/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {

	RenderTable(app)
	if err := WatchEdges(app); err != nil {
		return err
	}

	// Return success
	return nil
}
