package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.Logger
	gopi.TradfriManager
	gopi.ServiceDiscovery

	// Timeouts
	lookup *time.Duration
	id     *string
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SERVICE_TRADFRI        = "_coap._udp"
	SERVICE_LOOKUP_TIMEOUT = time.Second
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) Define(cfg gopi.Config) error {
	this.lookup = cfg.FlagDuration("mdns.timeout", SERVICE_LOOKUP_TIMEOUT, "Service lookup timeout")
	this.id = cfg.FlagString("tradfri.id", "", "Service identifier")
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.TradfriManager, this.ServiceDiscovery)
	return nil
}

func (this *app) Run(ctx context.Context) error {
	if id, srv, err := this.LookupService(ctx, *this.lookup); err != nil {
		return err
	} else if err := this.TradfriManager.Connect(id, srv.Host(), srv.Port()); err != nil {
		return err
	}

	// Get devices
	if devices, err := this.TradfriManager.Devices(ctx); err != nil {
		return err
	} else {
		this.Print("Observing devices")
		for _, device := range devices {
			ctx2, _ := context.WithCancel(ctx)
			go this.TradfriManager.ObserveDevice(ctx2, device)
		}
	}

	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()
	return nil
}

func (this *app) LookupService(parent context.Context, timeout time.Duration) (string, gopi.ServiceRecord, error) {
	// Lookup Gateway and return id and record for gateway connection
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	if gateways, err := this.ServiceDiscovery.Lookup(ctx, SERVICE_TRADFRI); err != nil {
		return "", nil, err
	} else if len(gateways) == 0 {
		return "", nil, gopi.ErrNotFound.WithPrefix("LookupService")
	} else if id := strings.TrimSpace(*this.id); id != "" {
		return id, gateways[0], nil
	} else {
		return gateways[0].Name(), gateways[0], nil
	}
}
