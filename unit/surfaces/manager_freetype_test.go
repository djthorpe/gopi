// +build freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces_test

import (
	"testing"

	// Frameworks

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/display"
	_ "github.com/djthorpe/gopi/v2/unit/fonts/freetype"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/platform"
	_ "github.com/djthorpe/gopi/v2/unit/surfaces"
)

func Test_Manager_Freetype_000(t *testing.T) {
	t.Log("Test_Manager_Freetype_000")
}
