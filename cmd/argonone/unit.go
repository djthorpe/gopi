package main

import (
	_ "github.com/djthorpe/gopi/v3/pkg/db/influxdb"  // Metric Writer
	_ "github.com/djthorpe/gopi/v3/pkg/dev/argonone" // Argon One
	_ "github.com/djthorpe/gopi/v3/pkg/hw/lirc"      // IR
	_ "github.com/djthorpe/gopi/v3/pkg/log"          // Logger
	_ "github.com/djthorpe/gopi/v3/pkg/mdns"         // RPC Service Discovery
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/client"   // RPC Client
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/ping"     // RPC Ping Service
)
