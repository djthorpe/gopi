// +build freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package freetype_test

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	"github.com/djthorpe/gopi/v2/unit/fonts"
	_ "github.com/djthorpe/gopi/v2/unit/fonts/freetype"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_Freetype_000(t *testing.T) {
	t.Log("Test_Freetype_000")
}

func Test_Freetype_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Freetype_001, []string{"-debug"}, "fonts"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Freetype_001(app gopi.App, t *testing.T) {
	fonts := app.Fonts()
	t.Log(fonts)
	if wd, err := os.Getwd(); err != nil {
		t.Error(err)
	} else {
		fontPath := filepath.Join(wd, "..", "..", "..", "etc", "fonts", "Damion", "Damion-Regular.ttf")
		if _, err := os.Stat(fontPath); os.IsNotExist(err) {
			t.Error(err)
		} else if face, err := fonts.OpenFace(fontPath); err != nil {
			t.Error(err)
		} else {
			t.Log(face)
		}
	}
}
func Test_Freetype_002(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Freetype_002, []string{"-debug"}, "fonts"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Freetype_002(app gopi.App, t *testing.T) {
	manager := app.Fonts()
	if wd, err := os.Getwd(); err != nil {
		t.Error(err)
	} else {
		fontPath := filepath.Join(wd, "..", "..", "..", "etc", "fonts", "Damion", "Damion-Regular.ttf")
		if _, err := os.Stat(fontPath); os.IsNotExist(err) {
			t.Error(err)
		} else if face, err := manager.OpenFace(fontPath); err != nil {
			t.Error(err)
		} else if image, err := face.(fonts.Face).BitmapForRunePixels('g', 512); err != nil {
			t.Error(err)
		} else if f, err := ioutil.TempFile("", "image*.png"); err != nil {
			t.Error(err)
		} else if err := png.Encode(f, image); err != nil {
			t.Error(err)
		} else {
			f.Close()
			fmt.Println(f.Name())
		}
	}
}
