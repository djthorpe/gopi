package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/log"
	"github.com/djthorpe/gopi/v3/pkg/platform"
)

type app struct {
	gopi.Unit
	*log.Log
	*platform.Platform
}

func (this *app) Run(ctx context.Context) error {
	fmt.Println(this.Platform)
	return nil
}

func (this *app) String() string {
	str := "<app"
	str += " platform=" + fmt.Sprint(this.Platform)
	str += " log=" + fmt.Sprint(this.Log)
	return str + ">"
}
