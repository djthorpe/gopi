package main

import (
	"context"
	"net/url"
	"os"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	table "github.com/djthorpe/gopi/v3/pkg/table"
)

func (this *app) RunCast(ctx context.Context, stub gopi.CastStub) error {
	casts, err := stub.ListCasts(ctx)
	if err != nil {
		return err
	} else if len(casts) == 0 {
		return gopi.ErrNotFound
	}

	// Display platform information
	table := table.New()

	table.SetHeader(header{"Name"}, header{"Id"}, header{"Model"}, header{"Service"}, header{"State"})
	for _, cast := range casts {
		table.Append(cast.Name(), cast.Id(), cast.Model(), cast.Service(), cast.State())
	}
	table.Render(os.Stdout)

	// Return success
	return nil
}

func (this *app) RunCastApp(ctx context.Context, stub gopi.CastStub) error {
	args := this.Args()
	if *this.castId == "" || len(args) != 1 {
		return gopi.ErrHelp
	}
	app := args[0]
	switch app {
	case "default":
		app = gopi.CAST_APPID_DEFAULT
	case "mutablemedia":
		app = gopi.CAST_APPID_MUTABLEMEDIA
	}
	if err := stub.SetApp(ctx, *this.castId, app); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *app) RunCastLoad(ctx context.Context, stub gopi.CastStub) error {
	args := this.Args()
	if *this.castId == "" || len(args) != 1 {
		return gopi.ErrHelp
	}
	if u, err := url.Parse(args[0]); err != nil {
		return err
	} else if err := stub.LoadURL(ctx, *this.castId, u); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *app) RunCastSeek(ctx context.Context, stub gopi.CastStub) error {
	args := this.Args()
	if *this.castId == "" || len(args) != 1 {
		return gopi.ErrHelp
	}
	if time, err := time.ParseDuration(args[0]); err != nil {
		return err
	} else if err := stub.SeekAbs(ctx, *this.castId, time); err != nil {
		return err
	}

	// Return success
	return nil
}
