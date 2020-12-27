// +build rpi

package rpi_test

import (
	"fmt"
	"testing"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v3/pkg/sys/rpi"
)

func Test_TVService_000(t *testing.T) {
	t.Log("Test_TVService_000")
}

func Test_TVService_001(t *testing.T) {
	if instance := rpi.VCHI_Init(); instance == nil {
		t.Error("VCHI_Init failed")
	} else if _, err := rpi.VCHI_TVInit(instance); err != nil {
		t.Error("VCHI_TVInit failed: ", err)
	} else if err := rpi.VCHI_TVStop(instance); err != nil {
		t.Error("VCHI_TVStop failed: ", err)
	}
}

func Test_TVService_002(t *testing.T) {
	if instance := rpi.VCHI_Init(); instance == nil {
		t.Error("VCHI_Init failed")
	} else if _, err := rpi.VCHI_TVInit(instance); err != nil {
		t.Error("VCHI_TVInit failed: ", err)
	} else {
		defer rpi.VCHI_TVStop(instance)
		if displays, err := rpi.VCHI_TVGetAttachedDevices(); err != nil {
			t.Error("VCHI_TVGetAttachedDevices failed: ", err)
		} else {
			for _, display := range displays {
				if state, err := rpi.VCHI_TVGetDisplayState(display); err != nil {
					t.Error("VCHI_TVGetDisplayState failed: ", err)
				} else {
					t.Log(display, state)
				}
			}
		}
	}
}

func Test_TVService_003(t *testing.T) {
	if instance := rpi.VCHI_Init(); instance == nil {
		t.Error("VCHI_Init failed")
	} else if _, err := rpi.VCHI_TVInit(instance); err != nil {
		t.Error("VCHI_TVInit failed: ", err)
	} else {
		defer rpi.VCHI_TVStop(instance)
		if displays, err := rpi.VCHI_TVGetAttachedDevices(); err != nil {
			t.Error("VCHI_TVGetAttachedDevices failed: ", err)
		} else {
			for _, display := range displays {
				if state, err := rpi.VCHI_TVGetDisplayState(display); err != nil {
					t.Error("VCHI_TVGetDisplayState failed: ", err)
				} else if info, err := rpi.VCHI_TVGetDisplayInfo(display); err != nil {
					t.Error("VCHI_TVGetDisplayInfo failed: ", err)
				} else {
					t.Log(display, state, info)
					if state.Flags()&rpi.TV_STATE_HDMI_ATTACHED > 0 {
						if err := rpi.VCHI_TVHDMIPowerOnPreferred(display); err != nil {
							t.Error("VCHI_TVHDMIPowerOnPreferred failed: ", err)
						}
						if err := rpi.VCHI_TVPowerOff(display); err != nil {
							t.Error("VCHI_TVHDMIPowerOff failed: ", err)
						}
					}
				}
			}
		}
	}
}

func Test_TVService_004(t *testing.T) {
	instance := rpi.VCHI_Init()
	if instance == nil {
		t.Fatal("VCHI_Init failed")
	} else if _, err := rpi.VCHI_TVInit(instance); err != nil {
		t.Fatal("VCHI_TVInit failed: ", err)
	}
	defer rpi.VCHI_TVStop(instance)
	rpi.VCTV_RegisterCallback(func(evt rpi.TVDisplayStateFlag, id rpi.DXDisplayId) {
		fmt.Println("Event", evt, "Display", id)
	})
	fmt.Println("Waiting for event")
	time.Sleep(60 * time.Second)
}
