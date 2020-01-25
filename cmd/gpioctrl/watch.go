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
	"strconv"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////

func pins(pins string) ([]gopi.GPIOPin, error) {
	logical_pins := make([]gopi.GPIOPin, 0)
	for _, value := range strings.Split(pins, ",") {
		if pin, err := strconv.ParseUint(value, 10, 64); err != nil {
			return nil, err
		} else {
			logical_pins = append(logical_pins, gopi.GPIOPin(pin))
		}
	}
	return logical_pins, nil
}

func WatchEdges(app gopi.App) error {
	gpio := app.GPIO()
	if edges, err := pins(app.Flags().GetString("edge", gopi.FLAG_NS_DEFAULT)); err != nil {
		return err
	} else if len(edges) == 0 {
		return nil
	} else {
		for _, logical := range edges {
			gpio.Watch(logical, gopi.GPIO_EDGE_BOTH)
		}
	}

	// Wait for CTRL+C
	fmt.Println("Press CTRL+C to end")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}
