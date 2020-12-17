package config_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/config"
)

func Test_Command_000(t *testing.T) {
	if cmd := config.NewCommand("name", "usage", "syntax", nil, nil); cmd == nil {
		t.Error("Unexpected nil return")
	} else {
		t.Log(cmd)
	}
}
