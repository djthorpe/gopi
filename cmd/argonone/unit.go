package main

import (
	_ "github.com/djthorpe/gopi/v3/pkg/db/influxdb"  // Metric Writer
	_ "github.com/djthorpe/gopi/v3/pkg/dev/argonone" // Argon One
	_ "github.com/djthorpe/gopi/v3/pkg/hw/lirc"      // IR
	_ "github.com/djthorpe/gopi/v3/pkg/log"          // Logger
)
