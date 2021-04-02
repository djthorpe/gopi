package main

import (
	_ "github.com/djthorpe/gopi/v3/pkg/event"      // gopi.Publisher
	_ "github.com/djthorpe/gopi/v3/pkg/log"        // gopi.Logger
	_ "github.com/djthorpe/gopi/v3/pkg/mdns"       // Multicast DNS
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/client" // gRPC Client
)

/*
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/metrics" // Metrics Service
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/ping"    // Ping Service
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/server"  // gRPC Server
	_ "github.com/djthorpe/gopi/v3/pkg/metrics"     // Metrics Platform
*/
