package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	// Modules

	"github.com/djthorpe/data"
	table "github.com/djthorpe/data/pkg/table"
	gopi "github.com/djthorpe/gopi/v3"

	// Dependencies
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/chromecast"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Chromecast struct {
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Chromecast) Define(cfg gopi.Config) {
	cfg.Command("cast", "Watch for Chromecast Events", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		ch := make(chan gopi.CastEvent)
		go func() {
			fmt.Println("Watching for events, press CTRL+C to end")
			for evt := range ch {
				fmt.Println(evt)
			}
		}()
		stub.Stream(ctx, ch)
		close(ch)
		return nil
	})

	cfg.Command("cast list", "Return list of Chromecasts", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		casts, err := stub.List(ctx, time.Second)
		if err != nil {
			return err
		}
		table := table.NewTable("Id", "Name", "Model", "Service", "State")
		for _, cast := range casts {
			table.Append(cast.Id(), cast.Name(), cast.Model(), cast.Service(), cast.State())
		}
		return table.Write(os.Stdout, table.OptHeader(), table.OptAscii(80, data.BorderLines))
	})

	cfg.Command("cast connect", "Connect to chromecast", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		args := this.GetArgs(ctx)
		if len(args) != 1 {
			return gopi.ErrBadParameter
		}
		if _, err := stub.Connect(ctx, args[0]); err != nil {
			return err
		} else {
			return nil
		}
	})

	cfg.Command("cast disconnect", "Disconnect from a chromecast", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		args := this.GetArgs(ctx)
		if len(args) != 1 {
			return gopi.ErrBadParameter
		}
		return stub.Disconnect(ctx, args[0])
	})

	cfg.Command("cast vol", "Set volume", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		args := this.GetArgs(ctx)
		if len(args) != 2 {
			return gopi.ErrBadParameter
		} else if vol, err := strconv.ParseFloat(args[1], 32); err != nil {
			return err
		} else if _, err := stub.SetVolume(ctx, args[0], float32(vol)); err != nil {
			return err
		}
		return nil
	})

	cfg.Command("cast mute", "Mute volume", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		args := this.GetArgs(ctx)
		if len(args) != 1 {
			return gopi.ErrBadParameter
		} else if _, err := stub.SetMuted(ctx, args[0], true); err != nil {
			return err
		}
		return nil
	})

	cfg.Command("cast unmute", "Unmute volume", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		args := this.GetArgs(ctx)
		if len(args) != 1 {
			return gopi.ErrBadParameter
		} else if _, err := stub.SetMuted(ctx, args[0], false); err != nil {
			return err
		}
		return nil
	})

	cfg.Command("cast app", "Launch application", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		args := this.GetArgs(ctx)
		if len(args) != 2 {
			return gopi.ErrBadParameter
		} else if _, err := stub.LaunchAppWithId(ctx, args[0], toAppId(args[1])); err != nil {
			return err
		}
		return nil
	})

}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *Chromecast) GetStub(ctx context.Context) gopi.CastStub {
	return ctx.Value(KeyStub).(gopi.CastStub)
}

func (this *Chromecast) GetArgs(ctx context.Context) []string {
	return ctx.Value(KeyArgs).([]string)
}

func toAppId(name string) string {
	switch name {
	case "backdrop":
		return gopi.CAST_APPID_BACKDROP
	case "mutablemedia":
		return gopi.CAST_APPID_MUTABLEMEDIA
	case "default":
		return gopi.CAST_APPID_DEFAULT
	default:
		return name
	}
}
