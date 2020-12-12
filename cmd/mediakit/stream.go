package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/djthorpe/gopi/v3"
)

func (this *app) Stream(ctx context.Context) error {
	args := this.Command.Args()
	if len(args) != 1 {
		return gopi.ErrHelp
	}

	if url, err := url.Parse(args[0]); err != nil {
		return err
	} else if media, err := this.MediaManager.OpenURL(url); err != nil {
		return err
	} else {
		defer this.MediaManager.Close(media)
		return this.StreamMedia(ctx, media)
	}
}

func (this *app) StreamMedia(ctx context.Context, media gopi.MediaInput) error {
	// Iterate through the frames decoding them
	return media.Read(ctx, nil, func(ctx gopi.MediaDecodeContext, packet gopi.MediaPacket) error {
		return media.DecodeFrameIterator(ctx, packet, func(frame gopi.MediaFrame) error {
			fmt.Println("f=", frame)
			return nil
		})
	})
}
