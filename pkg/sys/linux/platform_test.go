// +build linux

package linux_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

func Test_Platform_000(t *testing.T) {
	t.Log("Test_Platform_000")
}

func Test_Platform_001(t *testing.T) {
	if serial := linux.SerialNumber(); serial == "" {
		t.Error("Unexpected response from SerialNumber")
	} else {
		t.Log("serial", serial)
	}
}

func Test_Platform_002(t *testing.T) {
	if uptime := linux.Uptime(); uptime <= 0 {
		t.Error("Unexpected response from Uptime")
	} else {
		t.Log("uptime", uptime)
	}
}
func Test_Platform_003(t *testing.T) {
	if l1, l5, l15 := linux.LoadAverage(); l1 == 0 {
		t.Error("Unexpected response from LoadAverage")
	} else if l5 == 0 {
		t.Error("Unexpected response from LoadAverage")
	} else if l15 == 0 {
		t.Error("Unexpected response from LoadAverage")
	} else {
		t.Log("load averages", l1, l5, l15)
	}
}
