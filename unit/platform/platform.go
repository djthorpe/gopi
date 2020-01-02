/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package platform

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

type Platform struct{}

type platform struct {
	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Platform) Name() string { return "gopi.Platform" }

func (config Platform) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(platform)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *platform) String() string {
	return fmt.Sprintf("<gopi.Platform type=%v serial=%v uptime=%vhr>", this.Type(), strconv.Quote(this.SerialNumber()), this.Uptime().Truncate(time.Hour).Hours())
}
