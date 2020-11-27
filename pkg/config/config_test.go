package config_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/config"
)

func Test_Config_000(t *testing.T) {
	t.Log("Test_Config_000")
}

func Test_Config_001(t *testing.T) {
	if cfg := config.New(t.Name(), nil); cfg == nil {
		t.Error("Unexpected nil value")
	} else {
		t.Log(cfg)
	}
}

func Test_Config_002(t *testing.T) {
	if cfg := config.New(t.Name(), []string{"-debug"}); cfg == nil {
		t.Error("Unexpected nil value")
	} else if ptr := cfg.FlagBool("debug", false, "Test debug"); ptr == nil {
		t.Error("Unexpected nil value")
	} else if *ptr != false {
		t.Error("Unexpected default value")
	} else if err := cfg.Parse(); err != nil {
		t.Error(err)
	} else if *ptr != true {
		t.Error("Unexpected parsed value")
	} else if cfg.GetBool("debug") != true {
		t.Error("Unexpected GetBool value")
	}
}

func Test_Config_003(t *testing.T) {
	if cfg := config.New(t.Name(), []string{"-debug", "99"}); cfg == nil {
		t.Error("Unexpected nil value")
	} else if ptr := cfg.FlagUint("debug", 88, "Test debug"); ptr == nil {
		t.Error("Unexpected nil value")
	} else if *ptr != 88 {
		t.Error("Unexpected default value")
	} else if err := cfg.Parse(); err != nil {
		t.Error(err)
	} else if *ptr != 99 {
		t.Error("Unexpected parsed value")
	} else if cfg.GetUint("debug") != 99 {
		t.Error("Unexpected GetUint value")
	}
}

func Test_Config_004(t *testing.T) {
	if cfg := config.New(t.Name(), []string{"a", "b", "c"}); cfg == nil {
		t.Error("Unexpected nil value")
	} else if err := cfg.Parse(); err != nil {
		t.Error(err)
	} else if args := cfg.Args(); len(args) != 3 {
		t.Error("Unexpected Args() value:", args)
	} else {
		t.Log("args=", args)
	}
}
