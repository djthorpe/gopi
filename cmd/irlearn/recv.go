package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/djthorpe/gopi/v3"
)

var (
	hdr sync.Once
)

const (
	format = "%-10s %-20s %-15s %-15s\n"
)

func FormatDevice(device gopi.InputDevice, code uint32) string {
	device ^= gopi.INPUT_DEVICE_REMOTE
	str := strings.TrimPrefix(fmt.Sprint(device), "INPUT_DEVICE_")
	str += " " + FormatCode(code)
	return str
}

func FormatCode(code uint32) string {
	if code <= 0xFF {
		return fmt.Sprintf("0x%02X", code)
	} else if code <= 0xFFFF {
		return fmt.Sprintf("0x%04X", code)
	} else {
		return fmt.Sprintf("0x%08X", code)
	}
}

func FormatKey(key gopi.KeyCode) string {
	if key == gopi.KEYCODE_NONE {
		return "-"
	} else {
		return fmt.Sprint(key)
	}
}

func FormatType(t gopi.InputType) string {
	return strings.ToLower(strings.TrimPrefix(fmt.Sprint(t), "INPUT_EVENT_KEY"))
}

func FormatInputEvent(w io.Writer, evt gopi.InputEvent) {
	hdr.Do(func() {
		row := []interface{}{
			"EVENT", "NAME", "DEVICE", "KEY",
		}
		fmt.Fprintf(w, format, row...)
	})
	row := []interface{}{
		FormatType(evt.Type()),
		strconv.Quote(evt.Name()),
		FormatDevice(evt.Device()),
		FormatKey(evt.Key()),
	}
	fmt.Fprintf(w, format, row...)
}

func (this *app) Recv(ctx context.Context) error {
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	fmt.Println("Receiving IR Events, press CTRL+C to exit")

	for {
		select {
		case <-ctx.Done():
			return nil
		case evt := <-ch:
			if inputevent, ok := evt.(gopi.InputEvent); ok {
				FormatInputEvent(os.Stdout, inputevent)
			}
		}
	}
}
