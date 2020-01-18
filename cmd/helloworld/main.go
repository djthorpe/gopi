/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {

	// Print out name
	fmt.Println("Hello, " + app.Flags().GetString("name", gopi.FLAG_NS_DEFAULT))
	fmt.Println("Press CTRL+C to exit")

	// Wait for CTRL+C
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}
