/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/base"
)

type Config struct{}
type unit struct{ base.Unit }

func (Config) Name() string { return "gopi_test.Config" }

func (Config) New(log gopi.Logger) (gopi.Unit, error) {
	instance := new(unit)
	if err := instance.Init(log); err != nil {
		return nil, err
	} else {
		return instance, nil
	}
}

func Test_Gopi_000(t *testing.T) {
	t.Log("Test_Gopi_000")
}

func Test_Gopi_001(t *testing.T) {
	config := Config{}
	if unit, err := gopi.New(config, nil); err != nil {
		t.Error(err)
	} else if err := unit.Close(); err != nil {
		t.Error(err)
	} else {
		t.Log(unit)
	}
}

func Test_Gopi_002(t *testing.T) {
	NewPublisher := func() gopi.Publisher {
		return new(base.Publisher)
	}
	pubsub := NewPublisher()
	t.Log(pubsub)
}
