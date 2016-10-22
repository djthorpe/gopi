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
	"flag"
)

import (
	"../util"
)

var (
	flagLevel = flag.String("log","ANY","Logging level. Use ANY, NONE, FATAL, ERROR, WARN, INFO, DEBUG2, DEBUG")
)

func Function(log *util.LoggerDevice) {
	// do something here
	for i := 0; i < 10; i++ {
		log.Info("Counter = %v",i)
	}
}

func main() {
	flag.Parse()

	logger, err := util.Logger(util.StderrLogger{ })
	if err != nil {
		panic("Cannot open logger")
	}
	if err := logger.SetLevelFromString(*flagLevel); err != nil {
		panic(err)
	}

	logger.Info("Hello, %v","World")
	logger.Debug("Hello, again!")
	logger.Error("Oops, hello!")

	defer logger.Close()
}
