/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package main

import (
	"os"
	"fmt"
)

import (
	"../device/input"
)

func main() {
	driver, err := input.NewFT5406()
	if err != nil {
		fmt.Println("Error: ",err)
		os.Exit(-1)
	}
	defer driver.Close()

	fmt.Println("Device = ",driver)

	driver.ProcessEvents()
}
