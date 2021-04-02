package main

import (
	"context"
	"fmt"
	"strconv"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"

	// Dependencies
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/rotel"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Rotel struct{}

///////////////////////////////////////////////////////////////////////
// METHODS

func (this *Rotel) GetStub(ctx context.Context) gopi.RotelStub {
	return ctx.Value(KeyStub).(gopi.RotelStub)
}

func (this *Rotel) GetArgs(ctx context.Context) []string {
	return ctx.Value(KeyArgs).([]string)
}

func (this *Rotel) Define(cfg gopi.Config) {
	cfg.Command("rotel", "List Rotel state", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		ch := make(chan gopi.RotelEvent)
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
	cfg.Command("rotel off", "Power Off", func(ctx context.Context) error {
		return this.GetStub(ctx).SetPower(ctx, false)
	})
	cfg.Command("rotel on", "Power On", func(ctx context.Context) error {
		return this.GetStub(ctx).SetPower(ctx, true)
	})
	cfg.Command("rotel source", "Set Source", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 1 {
			return gopi.ErrBadParameter
		} else {
			return stub.SetSource(ctx, args[0])
		}
	})
	cfg.Command("rotel vol", "Set Volume (1-96)", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if vol, err := strconv.ParseUint(args[0], 10, 32); err != nil {
			return err
		} else {
			return stub.SetVolume(ctx, uint(vol))
		}
	})
	cfg.Command("rotel mute", "Mute Volume", func(ctx context.Context) error {
		return this.GetStub(ctx).SetMute(ctx, true)
	})
	cfg.Command("rotel bass", "Set Bass (-10 <> +10)", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if value, err := strconv.ParseInt(args[0], 10, 32); err != nil {
			return err
		} else {
			return stub.SetBass(ctx, int(value))
		}
	})
	cfg.Command("rotel treble", "Set Treble (-10 <> +10)", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if value, err := strconv.ParseInt(args[0], 10, 32); err != nil {
			return err
		} else {
			return stub.SetTreble(ctx, int(value))
		}
	})
	cfg.Command("rotel bypass", "Set Bypass (0,1)", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if value, err := strconv.ParseBool(args[0]); err != nil {
			return err
		} else {
			return stub.SetBypass(ctx, value)
		}
	})
	cfg.Command("rotel balance", "Set Balance (L,R,0)", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if len(args) != 1 {
			return gopi.ErrBadParameter
		} else {
			return stub.SetBalance(ctx, args[0])
		}
	})
	cfg.Command("rotel dimmer", "Set Dimmer (0,6)", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if value, err := strconv.ParseUint(args[0], 10, 32); err != nil {
			return err
		} else {
			return stub.SetDimmer(ctx, uint(value))
		}
	})
	cfg.Command("rotel play", "Send Play Command", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 0 {
			return gopi.ErrBadParameter
		} else {
			return stub.Play(ctx)
		}
	})
	cfg.Command("rotel stop", "Send Stop Command", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 0 {
			return gopi.ErrBadParameter
		} else {
			return stub.Stop(ctx)
		}
	})
	cfg.Command("rotel pause", "Send Pause Command", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 0 {
			return gopi.ErrBadParameter
		} else {
			return stub.Pause(ctx)
		}
	})
	cfg.Command("rotel next", "Send Next Track Command", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 0 {
			return gopi.ErrBadParameter
		} else {
			return stub.NextTrack(ctx)
		}
	})
	cfg.Command("rotel prev", "Send Previous Track Command", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		if args := this.GetArgs(ctx); len(args) != 0 {
			return gopi.ErrBadParameter
		} else {
			return stub.PrevTrack(ctx)
		}
	})
}
