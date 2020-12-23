// +build drm

package drm_test

import (
	"fmt"
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

func Test_DRM_000(t *testing.T) {
	devices := drm.Devices()
	if devices == nil {
		t.Skip("Skipping test, no devices available")
	}
	for _, device := range devices {
		for _, node := range device.Nodes() {
			t.Log(device, node)
		}
	}
}

func Test_DRM_001(t *testing.T) {
	devices := drm.Devices()
	if devices == nil {
		t.Skip("Skipping test, no devices available")
	}
	nodes := []string{}
	for _, device := range devices {
		if node := device.AvailableNode(drm.DRM_NODE_PRIMARY); node != "" {
			nodes = append(nodes, node)
		}
	}
	if len(nodes) == 0 {
		t.Skip("Skipping, no available primary node")
	} else {
		t.Log("Available Nodes=", nodes)
	}
}

func Test_DRM_002(t *testing.T) {
	devices := drm.Devices()
	if devices == nil {
		t.Skip("Skipping test, no devices available")
	}
	nodes := []string{}
	for _, device := range devices {
		if node := device.AvailableNode(drm.DRM_NODE_PRIMARY); node != "" {
			nodes = append(nodes, node)
		}
	}
	if len(nodes) == 0 {
		t.Skip("Skipping, no available primary node")
	}

	fh, err := drm.OpenPath(nodes[0])
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	res, err := drm.GetResources(fh.Fd())
	if err != nil {
		t.Fatal(err)
	}
	defer res.Free()
	t.Log(res)

	for _, plane := range drm.Planes(fh.Fd()) {
		t.Log("PLANE", plane)
		props := drm.GetPlaneProperties(fh.Fd(), plane)
		defer props.Free()
		for _, key := range props.Keys() {
			prop := drm.NewProperty(fh.Fd(), key)
			defer prop.Free()
			t.Log("   ", prop)
		}

		fmt.Println(props)
	}

}

/*
func Test_DRM_001(t *testing.T) {
	fh, err := drm.OpenDevice("card1")
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	if r, err := drm.GetResources(fh.Fd()); err != nil {
		t.Error(err)
	} else {
		defer r.Free()
		t.Log(r)
	}
}

func Test_DRM_002(t *testing.T) {
	fh, err := drm.OpenDevice("card1")
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	r, err := drm.GetResources(fh.Fd())
	if err != nil {
		t.Fatal(err)
	}
	defer r.Free()

	connectors := r.Connectors()
	if len(connectors) == 0 {
		t.Log("Skipping test as no connectors")
		t.SkipNow()
	}
	for _, id := range connectors {
		if connector, err := drm.GetConnector(fh.Fd(), id); err != nil {
			t.Error(err)
		} else {
			defer connector.Free()
			t.Log(connector)
		}
	}
}
*/
