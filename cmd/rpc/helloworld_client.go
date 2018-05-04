/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// An example RPC Client tool
package main

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc/grpc"
	_ "github.com/djthorpe/gopi/sys/rpc/mdns"

	// RPC Clients
	hw "github.com/djthorpe/gopi/rpc/grpc/helloworld"
	metrics "github.com/djthorpe/gopi/rpc/grpc/metrics"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	// Client Pool
	pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
	addr, _ := app.AppFlags.GetString("addr")

	// Lookup any service record for application
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	if records, err := pool.Lookup(ctx, "", addr, 1); err != nil {
		done <- gopi.DONE
		return err
	} else if len(records) == 0 {
		done <- gopi.DONE
		return gopi.ErrDeadlineExceeded
	} else if conn, err := pool.Connect(records[0], 0); err != nil {
		done <- gopi.DONE
		return err
	} else if err := RunHelloworld(app, conn); err != nil {
		done <- gopi.DONE
		return err
	} else if err := RunMetrics(app, conn); err != nil {
		done <- gopi.DONE
		return err
	} else if err := pool.Disconnect(conn); err != nil {
		done <- gopi.DONE
		return err
	}

	// Success
	done <- gopi.DONE
	return nil
}

func RunHelloworld(app *gopi.AppInstance, conn gopi.RPCClientConn) error {
	pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
	name, _ := app.AppFlags.GetString("name")
	if client_ := pool.NewClient("mutablelogic.Helloworld", conn); client_ == nil {
		return gopi.ErrAppError
	} else if client, ok := client_.(*hw.Client); ok == false {
		return gopi.ErrAppError
	} else if message, err := client.SayHello(name); err != nil {
		return err
	} else {
		fmt.Printf("%v says '%v'\n\n", conn.Name(), message)
		return nil
	}
}

func RunMetrics(app *gopi.AppInstance, conn gopi.RPCClientConn) error {
	pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
	if client_ := pool.NewClient("mutablelogic.Metrics", conn); client_ == nil {
		return gopi.ErrAppError
	} else if client, ok := client_.(*metrics.Client); ok == false {
		return gopi.ErrAppError
	} else if err := client.Ping(); err != nil {
		return err
	} else if metrics, err := client.HostMetrics(); err != nil {
		return err
	} else {
		fmt.Printf("Ping returned without error for connection %v\n", conn.Name())
		fmt.Printf("Metrics for remote host are: %v\n", metrics)
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/client/helloworld:grpc", "rpc/client/metrics:grpc")

	if cur, err := user.Current(); err == nil {
		config.AppFlags.FlagString("name", cur.Name, "Your name")
	} else {
		config.AppFlags.FlagString("name", "", "Your name")
	}
	config.AppFlags.FlagString("addr", "", "Gateway address")

	// Set the RPCServiceRecord for server discovery
	config.Service = "helloworld"

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main))
}
