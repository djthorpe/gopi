// +build drm

package drm_test

import (
	"testing"

	drm "github.com/djthorpe/gopi/v3/pkg/graphics/drm"
)

func Test_Resources_000(t *testing.T) {
	if fh, err := drm.OpenPrimaryDevice(); err != nil {
		t.Error(err)
	} else {
		defer fh.Close()
		t.Log(fh.Name())
	}
}

func Test_Resources_001(t *testing.T) {
	fh, err := drm.OpenPrimaryDevice()
	if err != nil {
		t.Skip("Skipping tests, missing primary device")
	}
	defer fh.Close()
	if res, err := drm.NewResources(fh.Fd()); err != nil {
		t.Error(err)
	} else if err := res.Dispose(); err != nil {
		t.Error(err)
	}
}

func Test_Resources_002(t *testing.T) {
	fh, err := drm.OpenPrimaryDevice()
	if err != nil {
		t.Skip("Skipping tests, missing primary device")
	}
	defer fh.Close()
	res, err := drm.NewResources(fh.Fd())
	if err != nil {
		t.Error(err)
	}
	defer res.Dispose()
	t.Log(res)
	connectors := res.NewActiveConnectors()
	if len(connectors) == 0 {
		t.Skip("Skipping tests, missing active connectors")
	}
	names := []string{}
	for _, connector := range connectors {
		for _, mode := range connector.Modes("", 0, false) {
			names = append(names, mode.Name())
		}
	}
	// Now iterate through to get modes again
	for _, name := range names {
		modes := []*drm.Mode{}
		for _, connector := range connectors {
			modes = append(modes, connector.Modes(name, 0, false)...)
		}
		t.Log(name, "=>", modes)
	}
	// Dispose connectors
	for _, connector := range connectors {
		if err := connector.Dispose(); err != nil {
			t.Error(err)
		}
	}
}

func Test_Resources_003(t *testing.T) {
	fh, err := drm.OpenPrimaryDevice()
	if err != nil {
		t.Skip("Skipping tests, missing primary device")
	}
	defer fh.Close()
	res, err := drm.NewResources(fh.Fd())
	if err != nil {
		t.Error(err)
	}
	defer res.Dispose()
	t.Log(res)
	connectors := res.NewActiveConnectors()
	if len(connectors) == 0 {
		t.Skip("Skipping tests, missing active connectors")
	}
	names := []string{}
	for _, connector := range connectors {
		for _, mode := range connector.Modes("", 0, false) {
			names = append(names, mode.Name())
		}
	}
	// Now iterate through to get modes again
	for _, name := range names {
		if conns, err := res.NewActiveConnectorsForMode(name, 0); err != nil {
			t.Error(err)
		} else if len(conns) == 0 {
			t.Error("Expected connector to be returned")
		} else {
			// Dispose connectors
			for _, conn := range conns {
				if err := conn.Dispose(); err != nil {
					t.Error(err)
				}
			}
		}
	}
	// Dispose connectors
	for _, connector := range connectors {
		if err := connector.Dispose(); err != nil {
			t.Error(err)
		}
	}
}

func Test_Resources_004(t *testing.T) {
	fh, err := drm.OpenPrimaryDevice()
	if err != nil {
		t.Skip("Skipping tests, missing primary device")
	}
	defer fh.Close()
	res, err := drm.NewResources(fh.Fd())
	if err != nil {
		t.Error(err)
	}
	defer res.Dispose()
	connectors := res.NewActiveConnectors()
	if len(connectors) == 0 {
		t.Skip("Skipping tests, missing active connectors")
	}
	modes := []*drm.Mode{}
	for _, connector := range connectors {
		for _, mode := range connector.Modes("", 0, true) {
			modes = append(modes, mode)
		}
		if err := connector.Dispose(); err != nil {
			t.Error(err)
		}
	}
	t.Log("Preferred modes=", modes)
}

func Test_Resources_005(t *testing.T) {
	fh, err := drm.OpenPrimaryDevice()
	if err != nil {
		t.Skip("Skipping tests, missing primary device")
	}
	defer fh.Close()
	res, err := drm.NewResources(fh.Fd())
	if err != nil {
		t.Error(err)
	}
	defer res.Dispose()
	connectors := res.NewActiveConnectors()
	if len(connectors) == 0 {
		t.Skip("Skipping tests, missing active connectors")
	}
	for _, connector := range connectors {
		if mode := connector.PreferredMode("", 0); mode == nil {
			t.Error("Expected preferred mode to be returned")
		} else if encoder, err := res.NewEncoderForConnector(connector); err != nil {
			t.Error(err)
		} else if encoder == nil {
			t.Error("Expected encoder to be returned")
		} else if crtc, err := res.NewCrtcForEncoder(encoder); err != nil {
			t.Error(err)
		} else if crtc == nil {
			t.Error("Expected crtc to be returned")
		} else {
			t.Log("Preferred mode=", mode)
			t.Log("Encoder=", encoder)
			t.Log("Crtc=", crtc)
			if err := crtc.Dispose(); err != nil {
				t.Error(err)
			}
			if err := encoder.Dispose(); err != nil {
				t.Error(err)
			}
		}
		if err := connector.Dispose(); err != nil {
			t.Error(err)
		}
	}
}
