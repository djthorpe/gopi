package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/table"
)

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *app) RunCast(ctx context.Context) error {
	devices, err := this.GetCastDevices(ctx, this.Args())
	if err != nil {
		return err
	}

	table := table.New()
	table.SetHeader("Name", "Model", "ID", "Service", "State")
	for _, device := range devices {
		table.Append(header{device.Name()}, device.Model(), device.Id(), device.Service(), device.State())
	}
	table.Render(os.Stdout)

	// If there is one device then connect
	if len(devices) == 1 {
		return this.RunCastDevice(ctx, devices[0])
	}

	// Return success
	return nil
}

func (this *app) RunCastDevice(ctx context.Context, device gopi.Cast) error {
	if err := this.CastManager.Connect(device); err != nil {
		return err
	}

	fmt.Println("Device=", device)
	fmt.Println("Waiting for CTRL+C")
	<-ctx.Done()

	return this.CastManager.Disconnect(device)
}

func (this *app) GetCastDevices(ctx context.Context, filter []string) ([]gopi.Cast, error) {
	if this.CastManager == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("CastManager")
	}

	// Discover devices
	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	devices, err := this.CastManager.Devices(ctx2)
	if err != nil {
		return nil, err
	} else if ctx2.Err() != context.DeadlineExceeded {
		return nil, ctx2.Err()
	}

	// Where there are no filters, return all devices
	if len(filter) == 0 {
		return devices, nil
	}

	// Filter devices
	castmap := make(map[string]gopi.Cast, len(devices))
	result := make([]gopi.Cast, 0, len(devices))
	for _, device := range devices {
		key := device.Id()
		castmap[key] = device
	}
	for _, key := range filter {
		if device, exists := castmap[key]; exists {
			result = append(result, device)
		} else {
			return nil, gopi.ErrNotFound.WithPrefix(key)
		}
	}
	return result, nil
}
