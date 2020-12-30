package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/table"
)

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *app) RunCast(ctx context.Context) error {
	devices, err := this.GetCastDevices(ctx)
	if err != nil {
		return err
	} else if len(devices) == 0 {
		return gopi.ErrNotFound.WithPrefix("Cast")
	}

	table := table.New()
	table.SetHeader("Name", "Model", "ID", "Service", "State")
	for _, device := range devices {
		table.Append(header{device.Name()}, device.Model(), device.Id(), device.Service(), device.State())
	}
	table.Render(os.Stdout)

	// Return success
	return nil
}

func (this *app) RunCastApp(ctx context.Context) error {
	devices, err := this.GetCastDevices(ctx)
	if err != nil {
		return err
	} else if len(devices) != 1 {
		return gopi.ErrNotFound.WithPrefix("Cast")
	}

	// Watch in background
	if *this.watch {
		this.WaitGroup.Add(1)
		go this.RunCastWatch(ctx)
	}

	if args := this.Args(); len(args) != 1 {
		return gopi.ErrBadParameter.WithPrefix("AppId required")
	} else if err := this.CastManager.LaunchAppWithId(devices[0], args[0]); err != nil {
		return err
	}

	// Wait for watching to end
	this.WaitGroup.Wait()

	// Return success
	return nil
}

func (this *app) RunCastVol(ctx context.Context) error {
	devices, err := this.GetCastDevices(ctx)
	if err != nil {
		return err
	} else if len(devices) != 1 {
		return gopi.ErrNotFound.WithPrefix("Cast")
	}

	// Watch in background
	if *this.watch {
		this.WaitGroup.Add(1)
		go this.RunCastWatch(ctx)
	}

	if args := this.Args(); len(args) != 1 {
		return gopi.ErrBadParameter.WithPrefix("Volume required between 0.0 and 1.0")
	} else if value, err := strconv.ParseFloat(args[0], 32); err != nil {
		return err
	} else if err := this.CastManager.SetVolume(devices[0], float32(value)); err != nil {
		return err
	}

	// Wait for watching to end
	this.WaitGroup.Wait()

	// Return success
	return nil
}

func (this *app) RunCastLoad(ctx context.Context) error {
	devices, err := this.GetCastDevices(ctx)
	if err != nil {
		return err
	} else if len(devices) != 1 {
		return gopi.ErrNotFound.WithPrefix("Cast")
	}

	// Watch in background
	if *this.watch {
		this.WaitGroup.Add(1)
		go this.RunCastWatch(ctx)
	}

	// Connect and wait
	if err := this.CastManager.Connect(devices[0]); err != nil {
		return err
	}

	// Wait for app
	time.Sleep(time.Second)

	if args := this.Args(); len(args) != 1 {
		return gopi.ErrBadParameter.WithPrefix("Missing URL")
	} else if url, err := url.Parse(args[0]); err != nil {
		return err
	} else if err := this.CastManager.LoadURL(devices[0], url, true); err != nil {
		return err
	}

	// Wait for watching to end
	this.WaitGroup.Wait()

	// Return success
	return nil
}

func (this *app) RunCastMute(ctx context.Context) error {
	return this.RunCastMuteEx(ctx, true)
}

func (this *app) RunCastUnmute(ctx context.Context) error {
	return this.RunCastMuteEx(ctx, false)
}

func (this *app) RunCastMuteEx(ctx context.Context, value bool) error {
	devices, err := this.GetCastDevices(ctx)
	if err != nil {
		return err
	} else if len(devices) != 1 {
		return gopi.ErrNotFound.WithPrefix("Cast")
	}

	// Watch in background
	if *this.watch {
		this.WaitGroup.Add(1)
		go this.RunCastWatch(ctx)
	}

	if err := this.CastManager.SetMuted(devices[0], value); err != nil {
		return err
	}

	// Wait for watching to end
	this.WaitGroup.Wait()

	// Return success
	return nil
}

func (this *app) RunCastWatch(ctx context.Context) {
	defer this.WaitGroup.Done()
	fmt.Println("Press CTRL+C to end")

	evts := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(evts)

FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		case evt := <-evts:
			if castevt, ok := evt.(gopi.CastEvent); ok {
				this.Print(castevt)
			}
		}
	}
}

func (this *app) GetCastDevices(ctx context.Context) ([]gopi.Cast, error) {
	if this.CastManager == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("CastManager")
	}

	// Discover devices
	ctx2, cancel := context.WithTimeout(ctx, *this.timeout)
	defer cancel()
	devices, err := this.CastManager.Devices(ctx2)
	if err != nil {
		return nil, err
	} else if ctx2.Err() != context.DeadlineExceeded {
		return nil, ctx2.Err()
	}

	// Where there are no filters, return all devices
	if *this.name == "" {
		return devices, nil
	}

	// Filter devices
	castmap := make(map[string]gopi.Cast, len(devices))
	for _, device := range devices {
		castmap[device.Id()] = device
		castmap[device.Name()] = device
	}
	if device, exists := castmap[*this.name]; exists {
		return []gopi.Cast{device}, nil
	} else {
		return []gopi.Cast{}, nil
	}
}
