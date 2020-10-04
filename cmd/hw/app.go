package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/log"
	"github.com/djthorpe/gopi/v3/pkg/platform"
	"github.com/djthorpe/gopi/v3/pkg/spi"
)

type app struct {
	gopi.Unit
	*log.Log
	*platform.Platform
	*spi.Spi
}

func (this *app) Run(ctx context.Context) error {
	fmt.Println(this.Platform)
	fmt.Println(this.Spi)
	return nil
}

func (this *app) String() string {
	str := "<app"
	str += " platform=" + fmt.Sprint(this.Platform)
	str += " log=" + fmt.Sprint(this.Log)
	return str + ">"
}
