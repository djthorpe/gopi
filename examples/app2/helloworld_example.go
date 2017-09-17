/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// The canonical hello world example
package main

import (
	"fmt"

	gopi "github.com/djthorpe/gopi"
)

func main() {
	config, err := gopi.NewAppConfigFromJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
