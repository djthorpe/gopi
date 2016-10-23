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
	flagLevel = flag.String("level","ANY","Logging level. Use ANY, NONE, FATAL, ERROR, WARN, INFO, DEBUG, DEBUG2")
	flagFile = flag.String("file","","Logging file. If empty, logs to stderr")
	flagAppend = flag.Bool("append",false,"Append output to file")
)

func Function(log *util.LoggerDevice) {
	// do something here
	for i := 0; i < 10; i++ {
		log.Info("Counter = %v",i)
	}
}

func main() {
	flag.Parse()

	var logger *util.LoggerDevice
	var err error

	// Open logger
	if len(*flagFile) != 0 {
		logger, err = util.Logger(util.FileLogger{ Filename: *flagFile, Append: *flagAppend })
	} else {
		logger, err = util.Logger(util.StderrLogger{ })
	}
	if err != nil {
		panic("Cannot open logger: " + err.Error())
	}
	defer logger.Close()

	// Set logging level
	if err := logger.SetLevelFromString(*flagLevel); err != nil {
		panic("Cannot set log level: " + err.Error())
	}

	// Generate log messages
	logger.Info("Hello, %v","World")
	logger.Debug("Hello, again!")
	logger.Error("Oops, hello!")

	// Call function
	Function(logger)
}
