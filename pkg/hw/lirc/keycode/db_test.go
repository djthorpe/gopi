package keycode_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/hw/lirc/keycode"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type DatabaseApp struct {
	gopi.Unit
	*keycode.Cache
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Database_001(t *testing.T) {
	tool.Test(t, nil, new(DatabaseApp), func(app *DatabaseApp) {
		if app.Cache == nil {
			t.Error("nil Cache unit")
		}
	})
}

func Test_Database_002(t *testing.T) {
	tool.Test(t, nil, new(DatabaseApp), func(app *DatabaseApp) {
		if app.Cache == nil {
			t.Error("nil Cache unit")
		}
		if db, err := keycode.NewDatabase("../../../../etc/keycode/test.keycode", "test", app.Cache); err != nil {
			t.Error(err)
		} else if k := db.Lookup(gopi.INPUT_DEVICE_RC5_14, 1); k != gopi.KEYCODE_1 {
			t.Error("Unepected lookup for ", gopi.INPUT_DEVICE_RC5_14, 1, ":", k)
		} else if k := db.Lookup(gopi.INPUT_DEVICE_RC5_14, 2); k != gopi.KEYCODE_2 {
			t.Error("Unepected lookup for ", gopi.INPUT_DEVICE_RC5_14, 1, ":", k)
		} else if k := db.Lookup(gopi.INPUT_DEVICE_RC5_14, 3); k != gopi.KEYCODE_3 {
			t.Error("Unepected lookup for ", gopi.INPUT_DEVICE_RC5_14, 1, ":", k)
		} else if k := db.Lookup(gopi.INPUT_DEVICE_RC5_14, 4); k != gopi.KEYCODE_4 {
			t.Error("Unepected lookup for ", gopi.INPUT_DEVICE_RC5_14, 1, ":", k)
		} else {
			t.Log(db)
		}
	})
}

func Test_Database_003(t *testing.T) {
	tool.Test(t, nil, new(DatabaseApp), func(app *DatabaseApp) {
		if app.Cache == nil {
			t.Error("nil Cache unit")
		}
		if db, err := keycode.NewDatabase("../../../../etc/keycode/test.keycode", "test", app.Cache); err != nil {
			t.Error(err)
		} else if err := db.Set(gopi.INPUT_DEVICE_RC5_14, 0xFFFF, gopi.KEYCODE_BTN0); err != nil {
			t.Error(err)
		} else {
			t.Log(db)
		}
	})
}

func Test_Database_004(t *testing.T) {
	tool.Test(t, nil, new(DatabaseApp), func(app *DatabaseApp) {
		if app.Cache == nil {
			t.Error("nil Cache unit")
		}
		if db, err := keycode.NewDatabase("../../../../etc/keycode/test.keycode", "test", app.Cache); err != nil {
			t.Error(err)
		} else if err := db.Set(gopi.INPUT_DEVICE_RC5_14, 0xFFFF, gopi.KEYCODE_BTN0); err != nil {
			t.Error(err)
		} else if k := db.Lookup(gopi.INPUT_DEVICE_RC5_14, 0xFFFF); k != gopi.KEYCODE_BTN0 {
			t.Error("Unexpected lookup value after set")
		} else if db.Modified() == false {
			t.Error("Expected database dirty flag set")
		} else {
			t.Log(db)
		}
	})
}

func Test_Database_005(t *testing.T) {
	tool.Test(t, nil, new(DatabaseApp), func(app *DatabaseApp) {
		if app.Cache == nil {
			t.Error("nil Cache unit")
		}
		if db, err := keycode.NewDatabase("../../../../etc/keycode/test.keycode", "test", app.Cache); err != nil {
			t.Error(err)
		} else if err := db.Set(gopi.INPUT_DEVICE_RC5_14, 1, gopi.KEYCODE_BTN1); err != nil {
			t.Error(err)
		} else if k := db.Lookup(gopi.INPUT_DEVICE_RC5_14, 1); k != gopi.KEYCODE_BTN1 {
			t.Error("Unexpected lookup value after set", k)
		} else if db.Modified() == false {
			t.Error("Expected database dirty flag set")
		} else {
			t.Log(db)
		}
	})
}
