package main

import (
	_ "github.com/djthorpe/gopi/v3/pkg/dev/googlecast"
	_ "github.com/djthorpe/gopi/v3/pkg/event"
	_ "github.com/djthorpe/gopi/v3/pkg/graphics/fonts/freetype"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/display"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/gpio/broadcom"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/i2c"

	//_ "github.com/djthorpe/gopi/v3/pkg/hw/lirc"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/platform"
	_ "github.com/djthorpe/gopi/v3/pkg/log"
	_ "github.com/djthorpe/gopi/v3/pkg/mdns"
)
