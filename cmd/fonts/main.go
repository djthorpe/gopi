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
	"path/filepath"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
	tablewriter "github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func PrintFonts(app gopi.App) {
	table := tablewriter.NewWriter(os.Stdout)
	// Get filters
	family := app.Flags().GetString("family", gopi.FLAG_NS_DEFAULT)
	flags := gopi.FONT_FLAGS_NONE
	if app.Flags().GetBool("regular", gopi.FLAG_NS_DEFAULT) {
		flags |= gopi.FONT_FLAGS_STYLE_REGULAR
	}
	if app.Flags().GetBool("bold", gopi.FLAG_NS_DEFAULT) {
		flags |= gopi.FONT_FLAGS_STYLE_BOLD
	}
	if app.Flags().GetBool("italic", gopi.FLAG_NS_DEFAULT) {
		flags |= gopi.FONT_FLAGS_STYLE_ITALIC
	}
	if flags == gopi.FONT_FLAGS_NONE {
		flags = gopi.FONT_FLAGS_STYLE_ANY
	}

	// Output table
	table.SetHeader([]string{"Name (Index)", "Family", "Style", "Flags", "Glyphs"})
	for _, font := range app.Fonts().Faces(family, flags) {
		table.Append([]string{
			fmt.Sprintf("%v (%v)", font.Name(), font.Index()),
			font.Family(),
			font.Style(),
			strings.ToLower(strings.ReplaceAll(fmt.Sprint(font.Flags()), "FONT_FLAGS_", "")),
			fmt.Sprint(font.NumGlyphs()),
		})
	}

	table.Render()
}

func Main(app gopi.App, args []string) error {
	if len(args) > 0 {
		return gopi.ErrHelp
	}
	if fontPath := app.Flags().GetString("fonts.path", gopi.FLAG_NS_DEFAULT); fontPath == "" {
		return gopi.ErrBadParameter.WithPrefix("Missing -fonts.path flag")
	} else if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		return gopi.ErrBadParameter.WithPrefix("Invalid -fonts.path flag")
	} else if err := app.Fonts().OpenFacesAtPath(fontPath, func(_ gopi.FontManager, path string, info os.FileInfo) bool {
		// Ignore hidden files and folders
		if strings.HasPrefix(info.Name(), ".") {
			return false
		}
		// Recurse into folders
		if info.IsDir() {
			return true
		}
		// Ignore zero-sized files
		if info.Size() == 0 {
			return false
		}
		switch strings.ToLower(filepath.Ext(info.Name())) {
		case ".ttf", ".ttc", ".otf", ".otc":
			app.Log().Debug("Loading:", path)
			return true
		default:
			app.Log().Debug("Ignoring:", path)
			return false
		}
	}); err != nil {
		return err
	}

	// Print fonts
	PrintFonts(app)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, nil, "fonts"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Set flags
		app.Flags().FlagString("fonts.path", "", "Path to font library")
		app.Flags().FlagString("family", "", "Font family")
		app.Flags().FlagBool("regular", false, "Regular style filter")
		app.Flags().FlagBool("bold", false, "Bold style filter")
		app.Flags().FlagBool("italic", false, "Italic style filter")
		// Run and exit
		os.Exit(app.Run())
	}
}
