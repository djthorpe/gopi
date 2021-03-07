//+build openvg,rpi

package openvg_test

import (
	"testing"

	openvg "github.com/djthorpe/gopi/v3/pkg/sys/openvg"
	rpi "github.com/djthorpe/gopi/v3/pkg/sys/rpi"
)

func Test_OpenVG_001(t *testing.T) {
	if err := rpi.BCMHostInit(); err != nil {
		t.Error("Error:", err)
	}
	if vendor := openvg.GetString(openvg.VG_VENDOR); vendor == "" {
		t.Error("Unexpected nil return for VG_VENDOR")
	} else if renderer := openvg.GetString(openvg.VG_RENDERER); renderer == "" {
		t.Error("Unexpected nil return for VG_RENDERER")
	} else if version := openvg.GetString(openvg.VG_VERSION); version == "" {
		t.Error("Unexpected nil return for VG_VERSION")
	} else if ext := openvg.GetString(openvg.VG_EXTENSIONS); ext == "" {
		t.Error("Unexpected nil return for VG_EXTENSIONS")
	} else {
		t.Log("VG_VENDOR", vendor)
		t.Log("VG_RENDERER", renderer)
		t.Log("VG_VERSION", version)
		t.Log("VG_EXTENSIONS", ext)
	}
}
