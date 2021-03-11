package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.ConnPool
	gopi.Command
	gopi.Logger
	gopi.PingService
	gopi.MetricsService
	gopi.Server
	gopi.ServiceDiscovery
	gopi.Unit

	// Flags
	service, castId *string
	watch           *bool
}

func (this *app) Define(cfg gopi.Config) error {
	cfg.Command("server", "Start ping service", func(ctx context.Context) error {
		if network, addr, err := this.GetServeAddress(); err != nil {
			return err
		} else {
			return this.RunServer(ctx, network, addr)
		}
	})
	cfg.Command("version", "Display server version information", func(ctx context.Context) error {
		if stub, err := this.GetPingStub(); err != nil {
			return err
		} else {
			return this.RunVersion(ctx, stub)
		}
	})
	cfg.Command("ping", "Perform ping to server", func(ctx context.Context) error {
		if stub, err := this.GetPingStub(); err != nil {
			return err
		} else {
			return this.RunPing(ctx, stub)
		}
	})
	cfg.Command("metrics", "Retrieve metrics from server", func(ctx context.Context) error {
		if stub, err := this.GetMetricsStub(); err != nil {
			return err
		} else {
			return this.RunMetrics(ctx, stub)
		}
	})
	cfg.Command("cast", "List Google Chromecasts", func(ctx context.Context) error {
		if stub, err := this.GetGoogleCastStub(); err != nil {
			return err
		} else {
			return this.RunCast(ctx, stub)
		}
	})
	cfg.Command("cast app", "Start Chromecast Application", func(ctx context.Context) error {
		if stub, err := this.GetGoogleCastStub(); err != nil {
			return err
		} else {
			return this.RunCastApp(ctx, stub)
		}
	})
	cfg.Command("cast load", "Load media from URL", func(ctx context.Context) error {
		if stub, err := this.GetGoogleCastStub(); err != nil {
			return err
		} else {
			return this.RunCastLoad(ctx, stub)
		}
	})
	cfg.Command("cast seek", "Seek within playing media", func(ctx context.Context) error {
		if stub, err := this.GetGoogleCastStub(); err != nil {
			return err
		} else {
			return this.RunCastSeek(ctx, stub)
		}
	})
	cfg.Command("cast pause", "Pause media playback", func(ctx context.Context) error {
		if stub, err := this.GetGoogleCastStub(); err != nil {
			return err
		} else {
			return this.RunCastPause(ctx, stub)
		}
	})
	cfg.Command("cast vol", "Set volume", func(ctx context.Context) error {
		if stub, err := this.GetGoogleCastStub(); err != nil {
			return err
		} else {
			return this.RunCastVol(ctx, stub)
		}
	})

	// Rotel
	cfg.Command("rotel", "List Rotel state", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else {
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
		}
	})
	cfg.Command("rotel off", "Power Off", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else {
			return stub.SetPower(ctx, false)
		}
	})
	cfg.Command("rotel on", "Power On", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else {
			return stub.SetPower(ctx, true)
		}
	})
	cfg.Command("rotel source", "Set Source", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 1 {
			return gopi.ErrBadParameter
		} else {
			return stub.SetSource(ctx, args[0])
		}
	})
	cfg.Command("rotel vol", "Set Volume (1-96)", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if vol, err := strconv.ParseUint(args[0], 10, 32); err != nil {
			return err
		} else {
			return stub.SetVolume(ctx, uint(vol))
		}
	})
	cfg.Command("rotel mute", "Mute Volume", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else {
			return stub.SetMute(ctx, true)
		}
	})
	cfg.Command("rotel bass", "Set Bass (-10 <> +10)", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if value, err := strconv.ParseInt(args[0], 10, 32); err != nil {
			return err
		} else {
			return stub.SetBass(ctx, int(value))
		}
	})
	cfg.Command("rotel treble", "Set Treble (-10 <> +10)", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if value, err := strconv.ParseInt(args[0], 10, 32); err != nil {
			return err
		} else {
			return stub.SetTreble(ctx, int(value))
		}
	})
	cfg.Command("rotel bypass", "Set Bypass (0,1)", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if value, err := strconv.ParseBool(args[0]); err != nil {
			return err
		} else {
			return stub.SetBypass(ctx, value)
		}
	})
	cfg.Command("rotel balance", "Set Balance (L,R) (1 <> 15)", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) == 0 {
			return stub.SetBalance(ctx, "L", 0)
		} else if len(args) != 2 {
			return gopi.ErrBadParameter
		} else if value, err := strconv.ParseUint(args[1], 10, 32); err != nil {
			return err
		} else {
			return stub.SetBalance(ctx, args[0], uint(value))
		}
	})
	cfg.Command("rotel dimmer", "Set Dimmer (0,6)", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 1 {
			return gopi.ErrBadParameter
		} else if value, err := strconv.ParseUint(args[0], 10, 32); err != nil {
			return err
		} else {
			return stub.SetDimmer(ctx, uint(value))
		}
	})
	cfg.Command("rotel play", "Send Play Command", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 0 {
			return gopi.ErrBadParameter
		} else {
			return stub.Play(ctx)
		}
	})
	cfg.Command("rotel stop", "Send Stop Command", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 0 {
			return gopi.ErrBadParameter
		} else {
			return stub.Stop(ctx)
		}
	})
	cfg.Command("rotel pause", "Send Pause Command", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 0 {
			return gopi.ErrBadParameter
		} else {
			return stub.Stop(ctx)
		}
	})
	cfg.Command("rotel next", "Send Next Track Command", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 0 {
			return gopi.ErrBadParameter
		} else {
			return stub.NextTrack(ctx)
		}
	})
	cfg.Command("rotel prev", "Send Previous Track Command", func(ctx context.Context) error {
		if stub, err := this.GetRotelStub(); err != nil {
			return err
		} else if args := this.Command.Args(); len(args) != 0 {
			return gopi.ErrBadParameter
		} else {
			return stub.PrevTrack(ctx)
		}
	})

	// Global flags
	this.service = cfg.FlagString("srv", "", "name, service:name or host:port")

	// Set flags for cast functions
	this.castId = cfg.FlagString("id", "", "Chromecast Id", "cast", "cast app", "cast load", "cast seek", "cast pause", "cast vol")

	// Set watch flag
	this.watch = cfg.FlagBool("watch", false, "Watch for events", "cast")

	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if cmd, err := cfg.GetCommand(cfg.Args()); err != nil {
		return gopi.ErrHelp
	} else if cmd == nil {
		return gopi.ErrHelp
	} else {
		this.Command = cmd
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

func (this *app) GetStub(name string) (gopi.ServiceStub, error) {
	// Timeout for lookup after 500ms
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	service := "grpc"
	if *this.service != "" {
		if strings.Contains(*this.service, ":") == false {
			service = service + ":" + *this.service
		} else {
			service = *this.service
		}
	}

	if conn, err := this.ConnPool.ConnectService(ctx, "tcp", service, 0); err != nil {
		return nil, err
	} else if stub := conn.NewStub(name); stub == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("Cannot create stub: ", name)
	} else {
		return stub, nil
	}
}

func (this *app) GetPingStub() (gopi.PingStub, error) {
	if stub, err := this.GetStub("gopi.ping.Ping"); err != nil {
		return nil, err
	} else {
		return stub.(gopi.PingStub), nil
	}
}

func (this *app) GetMetricsStub() (gopi.MetricsStub, error) {
	if stub, err := this.GetStub("gopi.metrics.Metrics"); err != nil {
		return nil, err
	} else {
		return stub.(gopi.MetricsStub), nil
	}
}

func (this *app) GetGoogleCastStub() (gopi.CastStub, error) {
	if stub, err := this.GetStub("gopi.googlecast.Manager"); err != nil {
		return nil, err
	} else {
		return stub.(gopi.CastStub), nil
	}
}

func (this *app) GetRotelStub() (gopi.RotelStub, error) {
	if stub, err := this.GetStub("gopi.rotel.Manager"); err != nil {
		return nil, err
	} else {
		return stub.(gopi.RotelStub), nil
	}
}

func (this *app) GetServeAddress() (string, string, error) {
	var network, addr string

	args := this.Args()
	switch {
	case len(args) == 0:
		network = "tcp"
		addr = ":0"
	case len(args) == 1:
		if _, _, err := net.SplitHostPort(args[0]); err != nil {
			return "", "", err
		} else {
			network = "tcp"
			addr = args[0]
		}
	default:
		return "", "", gopi.ErrBadParameter.WithPrefix(this.Command.Name())
	}

	// Return success
	return network, addr, nil
}
