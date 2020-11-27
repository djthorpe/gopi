// +build darwin

package darwin_test

import (
	"testing"

	// Frameworks
	darwin "github.com/djthorpe/gopi/v3/pkg/sys/darwin"
)

func Test_Platform_000(t *testing.T) {
	t.Log("Test_Platform_000")
}

func Test_Platform_001(t *testing.T) {
	if serial := darwin.SerialNumber(); serial == "" {
		t.Error("Unexpected response from SerialNumber")
	} else {
		t.Log("serial", serial)
	}
}

func Test_Platform_002(t *testing.T) {
	if uptime := darwin.Uptime(); uptime <= 0 {
		t.Error("Unexpected response from Uptime")
	} else {
		t.Log("uptime", uptime)
	}
}

func Test_Platform_003(t *testing.T) {
	if l1, l5, l15 := darwin.LoadAverage(); l1 == 0 {
		t.Error("Unexpected response from LoadAverage")
	} else if l5 == 0 {
		t.Error("Unexpected response from LoadAverage")
	} else if l15 == 0 {
		t.Error("Unexpected response from LoadAverage")
	} else {
		t.Log("load averages", l1, l5, l15)
	}
}

func Test_Platform_004(t *testing.T) {
	if product := darwin.Product(); product == "" {
		t.Error("Unexpected response from Product")
	} else {
		t.Log("product", product)
	}
}

func Test_Platform_005(t *testing.T) {
	if cputype := darwin.CPUType(); cputype == 0 {
		t.Error("Unexpected response from CPUType")
	} else {
		t.Log("cputype", cputype)
	}
}

func Test_Platform_006(t *testing.T) {
	cpu64 := darwin.CPU64Bit()
	t.Log("cpu64", cpu64)
}
