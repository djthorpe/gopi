/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/djthorpe/gopi/v2"
)

func Test_Error_000(t *testing.T) {
	t.Log("Test_Error_000")
}

func Test_Error_001(t *testing.T) {
	for err := gopi.ErrNone; err <= gopi.ErrMax; err++ {
		str := err.Error()
		if strings.HasPrefix(str, "[") {
			t.Error("Missing Error String for code", uint(err))
		} else {
			t.Log(uint(err), "=>", err.Error())
		}
	}
}

func Test_Error_002(t *testing.T) {
	for err := gopi.ErrNone; err <= gopi.ErrMax; err++ {
		err_ := err.WithPrefix("Prefix")
		str := err_.Error()
		if errors.Is(err_, err) == false {
			t.Error("errors.Is failed for code", uint(err))
		} else if strings.HasPrefix(str, "Prefix") == false {
			t.Error("WithPrefix failed for code", uint(err))
		} else {
			t.Log(uint(err), "=>", err_.Error())
		}
	}
}

func Test_Error_003(t *testing.T) {
	error := gopi.NewCompoundError()
	if error.ErrorOrSelf() != nil {
		t.Error("Unexpected return from ErrorOrSelf")
	}
}

func Test_Error_004(t *testing.T) {
	err := gopi.NewCompoundError(gopi.ErrInternalAppError)
	if err.ErrorOrSelf() != gopi.ErrInternalAppError {
		t.Error("Unexpected return from ErrorOrSelf")
	} else if errors.Is(err.ErrorOrSelf(), gopi.ErrInternalAppError) == false {
		t.Error("Unexpected return from ErrorOrSelf")
	} else if err.Is(gopi.ErrInternalAppError) == false {
		t.Error("Unexpected return from Is")
	}
}

func Test_Error_005(t *testing.T) {
	err := gopi.NewCompoundError(gopi.ErrInternalAppError, gopi.ErrNotFound.WithPrefix("Prefix"))
	if err.ErrorOrSelf() != err {
		t.Error("Unexpected return from ErrorOrSelf")
	} else if err.Is(gopi.ErrInternalAppError) == false {
		t.Error("Unexpected return from Is")
	} else if err.Is(gopi.ErrNotFound) == false {
		t.Error("Unexpected return from Is")
	} else {
		t.Log(err)
	}
}

func Test_Error_006(t *testing.T) {
	err := gopi.NewCompoundError(nil)
	if err.Is(gopi.ErrNone) == false {
		t.Error("Unexpected return from Is")
	} else if err.Is(nil) == false {
		t.Error("Unexpected return from Is")
	} else if err.ErrorOrSelf() != nil {
		t.Error("Unexpected return from ErrorOrSelf")
	}
}

func Test_Error_007(t *testing.T) {
	err := gopi.NewCompoundError(gopi.ErrInternalAppError)
	if err.Is(gopi.ErrNone) == true {
		t.Error("Unexpected return from Is")
	}
}

func Test_Error_008(t *testing.T) {
	err := gopi.NewCompoundError()
	if err.Is(gopi.ErrNone) == false {
		t.Error("Unexpected return from Is")
	} else if str := err.Error(); str != gopi.ErrNone.Error() {
		t.Error("Unexpected return from Error")
	}
}

func Test_Error_009(t *testing.T) {
	err := gopi.NewCompoundError(gopi.ErrNotImplemented)
	if str := err.Error(); str != gopi.ErrNotImplemented.Error() {
		t.Error("Unexpected return from Error")
	}
}
