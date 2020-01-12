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
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
)

////////////////////////////////////////////////////////////////////////////////

func SetBackground(mgr gopi.SurfaceManager) error {
	return mgr.Do(func(gopi.SurfaceManager) error {
		if _, err := mgr.CreateBackground(gopi.SURFACE_FLAG_BITMAP, 1.0); err != nil {
			return err
		} else {
			return nil
		}
	})
}

func SetBitmap(mgr gopi.SurfaceManager, bitmap gopi.Bitmap, origin gopi.Point) error {
	return mgr.Do(func(gopi.SurfaceManager) error {
		bitmap.ClearToColor(gopi.ColorRed)
		if _, err := mgr.CreateSurfaceWithBitmap(bitmap, 0, 1.0, 0, origin, gopi.ZeroSize); err != nil {
			return err
		} else {
			return nil
		}
	})
}

func Main(app gopi.App, args []string) error {
	if len(args) > 0 {
		return gopi.ErrHelp
	}

	// Set Background
	if err := SetBackground(app.Surfaces()); err != nil {
		return err
	}

	// Put red bitmap on screen
	if bitmap, err := app.Surfaces().CreateBitmap(0, gopi.Size{100, 100}); err != nil {
		return err
	} else if err := SetBitmap(app.Surfaces(), bitmap, gopi.Point{100, 100}); err != nil {
		return err
	}

	// Wait for key press
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, nil, "surfaces"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Run and exit
		os.Exit(app.Run())
	}
}
