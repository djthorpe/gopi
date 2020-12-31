package main

import (
	"context"
	"net/url"
	"os"

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
	if err := stub.SetApp(ctx, *this.castId, gopi.CAST_APPID_MUTABLEMEDIA); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *app) RunCastLoad(ctx context.Context, stub gopi.CastStub) error {
	if u, err := url.Parse("http://aurl/"); err != nil {
		return err
	} else if err := stub.LoadURL(ctx, *this.castId, u); err != nil {
		return err
	}

	// Return success
	return nil
}
