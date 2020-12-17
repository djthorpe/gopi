package config_test

import (
	"errors"
	"testing"

	"github.com/djthorpe/gopi/v3"
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

func Test_Config_005(t *testing.T) {
	if cfg := config.New(t.Name(), []string{"a", "b", "c"}); cfg == nil {
		t.Error("Unexpected nil value")
	} else if err := cfg.Command("   ", "Empty Command", nil); err == nil {
		t.Error("Expected error return")
	} else {
		t.Log(err)
	}
}

func Test_Config_006(t *testing.T) {
	cfg := config.New(t.Name(), []string{"a", "b", "c"})
	if cfg == nil {
		t.Error("Unexpected nil value")
	} else if err := cfg.Command("a", "Command A", nil); err != nil {
		t.Error(err)
	} else if err := cfg.Command("a b", "Command A B", nil); err != nil {
		t.Error(err)
	} else if err := cfg.Command("a b c", "Command A B C", nil); err != nil {
		t.Error(err)
	} else if err := cfg.Command("a d c", "Command A D C", nil); errors.Is(err, gopi.ErrNotFound) == false {
		t.Error("Unexpected error result")
	} else if err := cfg.Command("a b", "Command A B", nil); errors.Is(err, gopi.ErrDuplicateEntry) == false {
		t.Error("Unexpected error result")
	}

	if cmd := cfg.GetCommand([]string{"a", "b", "c"}); cmd == nil {
		t.Error("Unexpected nil result")
	} else if cmd.Name() != "a b c" {
		t.Error("Wrong command returned", cmd)
	} else if len(cmd.Args()) != 0 {
		t.Error("Unexpected number of arguments returned", cmd)
	}

	if cmd := cfg.GetCommand([]string{"a", "b", "c", "1", "2", "3"}); cmd == nil {
		t.Error("Unexpected nil result")
	} else if cmd.Name() != "a b c" {
		t.Error("Wrong command returned", cmd)
	} else if len(cmd.Args()) != 3 {
		t.Error("Unexpected number of arguments returned", cmd)
	} else {
		t.Log(cmd)
	}
	if cmd := cfg.GetCommand([]string{"a", "b"}); cmd == nil {
		t.Error("Unexpected nil result")
	} else if cmd.Name() != "a b" {
		t.Error("Wrong command returned", cmd)
	} else if len(cmd.Args()) != 0 {
		t.Error("Unexpected number of arguments returned", cmd)
	}
	if cmd := cfg.GetCommand([]string{"a", "b", "1", "2", "3"}); cmd == nil {
		t.Error("Unexpected nil result")
	} else if cmd.Name() != "a b" {
		t.Error("Wrong command returned", cmd)
	} else if len(cmd.Args()) != 3 {
		t.Error("Unexpected number of arguments returned", cmd)
	} else {
		t.Log(cmd)
	}

	if cmd := cfg.GetCommand([]string{"a"}); cmd == nil {
		t.Error("Unexpected nil result")
	} else if cmd.Name() != "a" {
		t.Error("Wrong command returned", cmd)
	} else if len(cmd.Args()) != 0 {
		t.Error("Unexpected number of arguments returned", cmd)
	}

	if cmd := cfg.GetCommand([]string{"a", "1", "2", "3"}); cmd == nil {
		t.Error("Unexpected nil result")
	} else if cmd.Name() != "a" {
		t.Error("Wrong command returned", cmd)
	} else if len(cmd.Args()) != 3 {
		t.Error("Unexpected number of arguments returned", cmd)
	} else {
		t.Log(cmd)
	}

	if cmd := cfg.GetCommand([]string{}); cmd == nil {
		t.Error("Unexpected nil result")
	} else if cmd.Name() != "a" {
		t.Error("Wrong command returned", cmd)
	} else if len(cmd.Args()) != 0 {
		t.Error("Unexpected number of arguments returned", cmd)
	}

	if cmd := cfg.GetCommand([]string{"1", "2", "3"}); cmd == nil {
		t.Error("Unexpected nil result")
	} else if cmd.Name() != "a" {
		t.Error("Wrong command returned", cmd)
	} else if len(cmd.Args()) != 3 {
		t.Error("Unexpected number of arguments returned", cmd)
	}
}

func Test_Config_007(t *testing.T) {
	cfg := config.New(t.Name(), nil)
	cfg.Command("a", "Command A", nil)
	cfg.Command("a b", "Command A B", nil)
	cfg.Command("a b c", "Command A B C", nil)
	cfg.FlagBool("g", false, "Flag g")
	cfg.FlagBool("a", false, "Flag a", "a")
	cfg.FlagBool("b", false, "Flag b", "a b", "a")
	cfg.FlagBool("c", false, "Flag c", "a b c", "a b", "a")

	cfg.Usage("")
	cfg.Usage("a")
}
