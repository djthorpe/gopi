package main

import (
	_ "github.com/djthorpe/gopi/v3/pkg/log"
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/client" // gRPC Client
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/ping"   // Ping Service
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/server" // gRPC Server
)
