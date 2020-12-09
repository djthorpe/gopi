// +build mimic

package mimic_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/mimic"
)

func Test_Mimic_000(t *testing.T) {
	mimic.Init()
	mimic.Exit()
}

func Test_Mimic_001(t *testing.T) {
	mimic.Init()
	if val, err := mimic.SetVoiceList("."); err != nil {
		t.Error(err)
	} else {
		t.Log(val)
	}
	mimic.Exit()
}
