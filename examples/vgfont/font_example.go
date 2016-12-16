/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example shows how to load fonts from one or more font directories
// and then display a list of fonts loaded, with various information
package main

import (
	"flag"
	"fmt"
	"os"
)

import (
	app "github.com/djthorpe/gopi/app"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {
	app.Logger.Info("Device=%v", app.Device)
	app.Logger.Info("Display=%v", app.Display)
	app.Logger.Info("EGL=%v", app.EGL)
	app.Logger.Info("OpenVG=%v", app.OpenVG)
	app.Logger.Info("Fonts=%v", app.Fonts)

	// display the list of fonts based on criteria
	family, exists := app.FlagSet.GetString("family")
	if exists == false {
		families := app.Fonts.GetFamilies()
		if len(families) == 0 {
			return app.Logger.Error("No font families loaded, use -fontpath path to locate fonts")
		}
		for _, family := range families {
			fmt.Println(family)
		}
	} else {
		flags := khronos.VG_FONT_STYLE_ANY
		bold, _ := app.FlagSet.GetBool("bold")
		italic, _ := app.FlagSet.GetBool("italic")
		switch {
		case bold && italic:
			flags = khronos.VG_FONT_STYLE_BOLDITALIC
		case bold:
			flags = khronos.VG_FONT_STYLE_BOLD
		case italic:
			flags = khronos.VG_FONT_STYLE_ITALIC
		}
		faces := app.Fonts.GetFaces(family, flags)
		if len(faces) == 0 {
			return app.Logger.Error("No such family '%s'",family)
		}
		format := "%3s %-20s %-20s\n"
		fmt.Printf(format,"ID","Family","Style")
		fmt.Printf("--------------------------------------------\n")
		for _, face := range faces {
			fmt.Printf(format,fmt.Sprintf("%03d",face.GetIndex()),face.GetFamily(),face.GetStyle())
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_VGFONT)

	// Family, italic and bold
	config.FlagSet.FlagString("family", "", "List fonts for one particular family")
	config.FlagSet.FlagBool("bold", false, "List bold fonts")
	config.FlagSet.FlagBool("italic", false, "List italic fonts")

	// Create the application
	myapp, err := app.NewApp(config)
	if err == flag.ErrHelp {
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(MyRunLoop); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
