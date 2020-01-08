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
	"os"
	"path/filepath"
	"strings"
	"testing"

	// Frameworks
	ft "github.com/djthorpe/gopi/v2/sys/freetype"
)

func Test_Freetype_000(t *testing.T) {
	t.Log("Test_Freetype_000")
}

func Test_Freetype_001(t *testing.T) {
	for status := ft.FT_ERROR_MIN; status <= ft.FT_ERROR_MAX; status++ {
		status_error := fmt.Sprint(status.Error())
		if strings.HasPrefix(status_error, "FT_ERROR_") {
			t.Logf("%v => %s", int(status), status_error)
		} else {
			t.Logf("No status error for value: %v", status)
		}
	}
}

func Test_Freetype_002(t *testing.T) {
	if library, err := ft.FT_Init(); err != nil {
		t.Error(err)
	} else if err := ft.FT_Destroy(library); err != nil {
		t.Error(err)
	}
}

func Test_Freetype_003(t *testing.T) {
	if library, err := ft.FT_Init(); err != nil {
		t.Error(err)
	} else if major, minor, patch := ft.FT_Library_Version(library); major == 0 {
		t.Error("unexpected major version, ", major)
	} else {
		t.Logf("version={%v,%v,%v}", major, minor, patch)
		if err := ft.FT_Destroy(library); err != nil {
			t.Error(err)
		}
	}
}

func Test_Freetype_004(t *testing.T) {
	if library, err := ft.FT_Init(); err != nil {
		t.Error(err)
	} else {
		defer ft.FT_Destroy(library)
		if fontPath, err := os.Getwd(); err != nil {
			t.Error(err)
		} else {
			fontPath = filepath.Join(fontPath, "..", "..", "etc", "fonts", "Damion", "Damion-Regular.ttf")
			if _, err := os.Stat(fontPath); os.IsNotExist(err) {
				t.Error(fontPath, err)
			} else if _, err := ft.FT_NewFace(library, fontPath, 0); err != nil {
				t.Error(err)
			}
		}
	}
}

func Test_Freetype_005(t *testing.T) {
	if library, err := ft.FT_Init(); err != nil {
		t.Error(err)
	} else {
		defer ft.FT_Destroy(library)
		if fontPath, err := os.Getwd(); err != nil {
			t.Error(err)
		} else {
			fontPath = filepath.Join(fontPath, "..", "..", "etc", "fonts", "Damion", "Damion-Regular.ttf")
			if _, err := os.Stat(fontPath); os.IsNotExist(err) {
				t.Error(fontPath, err)
			} else if face, err := ft.FT_NewFace(library, fontPath, 0); err != nil {
				t.Error(err)
			} else if err := ft.FT_DoneFace(face); err != nil {
				t.Error(err)
			}
		}
	}
}

func Test_Freetype_006(t *testing.T) {
	if library, err := ft.FT_Init(); err != nil {
		t.Error(err)
	} else {
		defer ft.FT_Destroy(library)
		if fontPath, err := os.Getwd(); err != nil {
			t.Error(err)
		} else {
			fontPath = filepath.Join(fontPath, "..", "..", "etc", "fonts", "Damion", "Damion-Regular.ttf")
			if _, err := os.Stat(fontPath); os.IsNotExist(err) {
				t.Error(fontPath, err)
			} else if face, err := ft.FT_NewFace(library, fontPath, 0); err != nil {
				t.Error(err)
			} else if err := ft.FT_SelectCharmap(face, ft.FT_ENCODING_UNICODE); err != nil {
				t.Error(err)
			}
		}
	}
}
