/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Example command for discovery of RPC microservices using mDNS
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	tablewriter "github.com/olekukonko/tablewriter"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc/mdns"
)

////////////////////////////////////////////////////////////////////////////////

var (
	lock  sync.Mutex
	gctx  context.Context
	table *tablewriter.Table
)

func InitContext() {
	lock.Lock()
}

func SetContext(ctx context.Context) {
	gctx = ctx
	defer lock.Unlock()
}

func GetContext() context.Context {
	lock.Lock()
	defer lock.Unlock()
	return gctx
}

////////////////////////////////////////////////////////////////////////////////

func PrintHeader() {
	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "TYPE", "ADDR:PORT", "TTL", "TEXT"})
	table.SetAutoMergeCells(true)
}

func PrintRecord(s *gopi.RPCServiceRecord) {
	// Gather up the addresses
	addr := ""
	for _, ip4 := range s.IP4 {
		addr += fmt.Sprintf("%v:%v\n", ip4.String(), s.Port)
	}
	for _, ip6 := range s.IP6 {
		addr += fmt.Sprintf("%v:%v\n", ip6.String(), s.Port)
	}
	table.Append([]string{
		s.Name,
		s.Type,
		strings.TrimSpace(addr),
		fmt.Sprint(s.TTL),
		strings.Join(s.Text, "\n"),
	})
	table.Render()
	table.ClearRows()

	/*	table.SetHeader("NAME", "TYPE", "ADDR:PORT", "TTL", "TEXT")
		fmt.Printf("%s\n", )
		fmt.Printf("%-4s %-20s %-20s %-10v %-10s\n", "", s.Type, "", s.Port, s.TTL)
		for _, txt := range s.Text {
			fmt.Printf("%-53s %s\n", "", txt)
		}*/
}

func EventLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	var once sync.Once

	discovery := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery)

	// Subscribe to record discovery
	c := discovery.Subscribe()

FOR_LOOP:
	for {
		select {
		case evt := <-c:
			if rpc_evt, ok := evt.(gopi.RPCEvent); rpc_evt != nil && ok {
				once.Do(PrintHeader)
				PrintRecord(rpc_evt.ServiceRecord())
			}
		case <-done:
			break FOR_LOOP
		}
	}

	// Stop listening for events
	discovery.Unsubscribe(c)

	return nil
}

func BrowseLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	if service, _ := app.AppFlags.GetString("service"); service == "" {
		return errors.New("Missing or invalid -service parameter")
	} else if service_qualified, err := gopi.RPCServiceType(service, gopi.RPC_FLAG_NONE); err != nil {
		return err
	} else {
		app.Logger.Info("Browsing for '%v'", service_qualified)
		if discovery, ok := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery); discovery == nil || ok == false {
			return errors.New("Missing or invalid discovery service")
		} else if err := discovery.Browse(GetContext(), service_qualified); err != nil {
			return err
		}
	}

	// Wait for done
	_ = <-done
	return nil
}

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {
	// Set parameters
	timeout, _ := app.AppFlags.GetDuration("timeout")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	SetContext(ctx)

	// Wait until CTRL+C is pressed
	app.WaitForSignalOrTimeout(timeout)

	// Perform cancel
	cancel()

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the lirc instance
	config := gopi.NewAppConfig("mdns")

	// Set flags
	config.AppFlags.FlagDuration("timeout", 1*time.Second, "Browse timeout")
	config.AppFlags.FlagString("service", "", "Service")

	// Init
	InitContext()

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop, BrowseLoop, EventLoop))
}
