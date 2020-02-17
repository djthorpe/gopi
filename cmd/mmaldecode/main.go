/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
	mmal "github.com/djthorpe/gopi/v2/sys/mmal"
)

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Run and exit
		os.Exit(app.Run())
	}
}

////////////////////////////////////////////////////////////////////////////////

func CreateComponent(name string) (mmal.MMAL_ComponentHandle, error) {
	var component mmal.MMAL_ComponentHandle
	if err := mmal.MMALComponentCreate(name, &component); err != nil {
		return nil, err
	} else {
		return component, nil
	}
}

func DestroyComponent(component mmal.MMAL_ComponentHandle) error {
	return mmal.MMALComponentDestroy(component)
}

func Main(app gopi.App, args []string) error {
	if renderer, err := CreateComponent(MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER); err != nil {
		return err
	} else {
		defer DestroyComponent(renderer)
	}
	if decoder, err := CreateComponent(MMAL_COMPONENT_DEFAULT_VIDEO_DECODER); err != nil {
		return err
	} else {
		defer DestroyComponent(renderer)
	}

	// Return success
	return nil
}
