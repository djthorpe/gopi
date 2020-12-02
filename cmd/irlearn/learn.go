package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/djthorpe/gopi/v3"
)

func (this *app) Learn(ctx context.Context) error {
	if len(this.Command.Args()) == 0 {
		return fmt.Errorf("Missing argument")
	}
	for _, name := range this.Command.Args() {
		codes := this.LIRCKeycodeManager.Keycode(name)
		if len(codes) == 0 {
			fmt.Println("No keys found for", name)
			continue
		}
		for i, code := range codes {
			fmt.Println(i+1, "of", len(codes))
			if err := this.LearnCode(ctx, code, *this.name); err != nil {
				return err
			}
			fmt.Println()
		}
	}
	return nil
}

func (this *app) LearnCode(ctx context.Context, code gopi.KeyCode, name string) error {
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	timer := time.NewTimer(5 * time.Second)
	fmt.Println("Learning", code, "for", strconv.Quote(name))
	fmt.Println("  ...press the key on the remote corresponding to", code, "or CTRL+C to end")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return nil
		case <-ch:
			//fmt.Println(evt)
		}
	}
}
