/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func Background2(app *gopi.AppInstance, start chan<- struct{}, stop <-chan struct{}) error {
	app.Logger.Info("Waiting on Background2...")
	time.Sleep(time.Second)
	// Send start signal, to show initialization has happened
	start <- gopi.DONE
	app.Logger.Info("Started Background2...")

FOR_LOOP:
	for {
		select {
		case <-stop:
			break FOR_LOOP
		}
	}

	app.Logger.Info("Finished Background2...")
	return fmt.Errorf("Error from Background 2")

}

func Background1(app *gopi.AppInstance, start chan<- struct{}, stop <-chan struct{}) error {
	// Send start signal, to show initialization has happened
	start <- gopi.DONE
	app.Logger.Info("Started Background1...")

FOR_LOOP:
	for {
		select {
		case <-stop:
			break FOR_LOOP
		}
	}

	app.Logger.Info("Finished Background1...")
	return fmt.Errorf("Error from Background 1")
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	app.Logger.Info("Started Main...")
	time.Sleep(time.Second * 2)

	app.Logger.Info("Sending DONE signal from main...")
	done <- gopi.DONE

	app.Logger.Info("...Finished Main")
	return gopi.ErrHelp
}

////////////////////////////////////////////////////////////////////////////////

var (
	GitTag, GitBranch, GitHash string
	GoBuildTime                string
)

func main() {
	// Create the configuration
	config := gopi.NewAppConfig()

	// Set version parameters
	config.AppFlags.SetParam(gopi.PARAM_GOVERSION, runtime.Version())
	config.AppFlags.SetParam(gopi.PARAM_GITTAG, GitTag)
	config.AppFlags.SetParam(gopi.PARAM_GITBRANCH, GitBranch)
	config.AppFlags.SetParam(gopi.PARAM_GITHASH, GitHash)
	config.AppFlags.SetParam(gopi.PARAM_GOBUILDTIME, GoBuildTime)

	// Run the command line tool
	os.Exit(gopi.CommandLineTool2(config, Main, Background1, Background2))
}
