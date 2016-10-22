/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example shows you how the logging works for the gopi package, using
// the concrete VideoCore logging mechanism as a logging output. You can
// use the standard logger, which outputs logging information to stderr by
// using:
//
//    logger := util.Logger(util.StderrLogger{ })
//    logger.SetLevel(LOG_ANY)
//
package main

import (
	"../util"
)

func Function(log *util.LoggerDevice) {
	// do something here
	for i := 0; i < 10; i++ {
		log.Info("Counter = %v",i)
	}
}

func main() {
	logger, err := util.Logger(util.StderrLogger{ })
	if err != nil {
		panic("Cannot open logger")
	}
	logger.SetLevel(util.LOG_ANY)

	logger.Info("Hello, %v","World")
	logger.Debug("Hello, again!")
	logger.Error("Oops, hello!")

	defer logger.Close()
}
