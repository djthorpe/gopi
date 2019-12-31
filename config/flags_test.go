/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package config_test

import (
	"errors"
	"testing"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/config"
)

func Test_Config_000(t *testing.T) {
	t.Log("Test_Config_000")
}

func Test_Config_001(t *testing.T) {
	flags := config.NewFlags("Test_Config_001")
	if flags.Name() != "Test_Config_001" {
		t.Error("Unexpected return value from Name()")
	}
	if flags.Parsed() != false {
		t.Error("Unexpected return value from Parsed()")
	}
	if len(flags.Args()) != 0 {
		t.Error("Unexpected return value from Args()")
	}
	if len(flags.Flags()) != 0 {
		t.Error("Unexpected return value from Flags()")
	}
}

func Test_Config_002(t *testing.T) {
	flags := config.NewFlags("Test_Config_002")
	if err := flags.Parse([]string{"-help"}); errors.Is(err, gopi.ErrHelp) == false {
		t.Error("Unexpected return value from Parse()", err)
	}
	if flags.Parsed() != true {
		t.Error("Unexpected return value from Parsed()")
	}
	if len(flags.Args()) != 0 {
		t.Error("Unexpected return value from Args()")
	}
}

func Test_Config_003(t *testing.T) {
	flags := config.NewFlags("Test_Config_003")
	if err := flags.Parse([]string{"a", "b", "c"}); err != nil {
		t.Error("Unexpected return value from Parse()", err)
	}
	if flags.Parsed() != true {
		t.Error("Unexpected return value from Parsed()")
	}
	if len(flags.Args()) != 3 {
		t.Error("Unexpected return value from Args()")
	}
}

func Test_Config_004(t *testing.T) {
	flags := config.NewFlags("Test_Config_004")
	flags.FlagBool("test", false, "Test Flag")
	flags.FlagBool("test2", true, "Test2 Flag")
	if err := flags.Parse([]string{"-test", "-test2=false", "b", "c"}); err != nil {
		t.Error("Unexpected return value from Parse()", err)
	}
	if flags.Parsed() != true {
		t.Error("Unexpected return value from Parsed()")
	}
	if len(flags.Args()) != 2 {
		t.Error("Unexpected return value from Args()")
	}
	if flags.Has("test") != true {
		t.Error("Unexpected return value from Has()")
	}
	if flags.Has("test2") != true {
		t.Error("Unexpected return value from Has()")
	}
	if flags.Has("test3") != false {
		t.Error("Unexpected return value from Has()")
	}
	if flags.GetBool("test") != true {
		t.Error("Unexpected return value from GetBool()")
	}
	if flags.GetBool("test2") != false {
		t.Error("Unexpected return value from GetBool()")
	}
	if flags.GetBool("test3") != false {
		t.Error("Unexpected return value from GetBool()")
	}
}

func Test_Config_005(t *testing.T) {
	flags := config.NewFlags("Test_Config_005")
	flags.FlagString("test", "test", "Test Flag")
	flags.FlagBool("test2", true, "Test2 Flag")
	flags.FlagString("test3", "default", "Test3 Flag")
	if err := flags.Parse([]string{"-test", "hello", "-test2=true", "b", "c"}); err != nil {
		t.Error("Unexpected return value from Parse()", err)
	}
	if flags.Parsed() != true {
		t.Error("Unexpected return value from Parsed()")
	}
	if len(flags.Args()) != 2 {
		t.Error("Unexpected return value from Args()")
	}
	if flags.GetString("test") != "hello" {
		t.Error("Unexpected return value from GetString()")
	}
	if flags.GetString("test2") != "true" {
		t.Error("Unexpected return value from GetString()")
	}
	if flags.GetString("test3") != "default" {
		t.Error("Unexpected return value from GetString()")
	}
}

func Test_Config_006(t *testing.T) {
	flags := config.NewFlags("Test_Config_005")
	flags.FlagInt("test", -1234, "Test Flag")
	flags.FlagUint("test2", 1234, "Test2 Flag")
	flags.FlagFloat64("test3", 12.34, "Test3 Flag")
	flags.FlagDuration("test4", 1234*time.Second, "Test4 Flag")
	if err := flags.Parse([]string{"a", "b"}); err != nil {
		t.Error("Unexpected return value from Parse()", err)
	}
	if flags.Parsed() != true {
		t.Error("Unexpected return value from Parsed()")
	}
	if len(flags.Args()) != 2 {
		t.Error("Unexpected return value from Args()")
	}
	if flags.GetString("test") != "-1234" {
		t.Error("Unexpected return value from GetString()")
	}
	if flags.GetInt("test") != -1234 {
		t.Error("Unexpected return value from GetInt()")
	}
	if flags.GetUint("test2") != 1234 {
		t.Error("Unexpected return value from GetUint()")
	}
	if flags.GetFloat64("test3") != 12.34 {
		t.Error("Unexpected return value from GetFloat64()")
	}
	if flags.GetDuration("test4") != 1234*time.Second {
		t.Error("Unexpected return value from GetDuration()")
	}
}
