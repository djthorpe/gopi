package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
	"github.com/olekukonko/tablewriter"

	_ "github.com/djthorpe/gopi/v3/pkg/graphics/fonts/freetype"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/display"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/gpio/broadcom"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/platform"
	_ "github.com/djthorpe/gopi/v3/pkg/log"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	rePin = regexp.MustCompile("^(GPIO|gpio)(\\d+)$")
)

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *app) RunGPIO(ctx context.Context) error {
	args := this.Args()
	if len(args) == 0 {
		return this.GPIOListPins(ctx, nil)
	}
	pins, err := this.GPIOParsePins(args[1:])
	if err != nil {
		return err
	} else if len(pins) == 0 {
		return gopi.ErrBadParameter
	}

	var result error
	switch args[0] {
	case "value":
		break
	case "pullup":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_INPUT)
			if err := this.GPIO.SetPullMode(pin, gopi.GPIO_PULL_UP); err != nil {
				result = multierror.Append(result, err)
			}
		}
	case "pulldown":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_INPUT)
			if err := this.GPIO.SetPullMode(pin, gopi.GPIO_PULL_DOWN); err != nil {
				result = multierror.Append(result, err)
			}
		}
	case "in", "input":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_INPUT)
		}
	case "out", "output":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_OUTPUT)
		}
	case "low", "0":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_OUTPUT)
			this.GPIO.WritePin(pin, gopi.GPIO_LOW)
		}
	case "high", "1":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_OUTPUT)
			this.GPIO.WritePin(pin, gopi.GPIO_HIGH)
		}
	case "alt0":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_ALT0)
		}
	case "alt1":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_ALT1)
		}
	case "alt2":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_ALT2)
		}
	case "alt3":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_ALT3)
		}
	case "alt4":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_ALT4)
		}
	case "alt5":
		for _, pin := range pins {
			this.GPIO.SetPinMode(pin, gopi.GPIO_ALT5)
		}
	default:
		return gopi.ErrBadParameter.WithPrefix(args[0])
	}

	return this.GPIOListPins(ctx, pins)
}

// List all pins in a table
func (this *app) GPIOListPins(ctx context.Context, pins []gopi.GPIOPin) error {
	// Check to make sure there are pins if individual pins are not specified
	if num := this.GPIO.NumberOfPhysicalPins(); len(pins) == 0 && num == 0 {
		return fmt.Errorf("No GPIO interface defined")
	} else if len(pins) == 0 {
		for pin := uint(1); pin <= num; pin++ {
			pins = append(pins, this.GPIO.PhysicalPin(pin))
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Physical", "Logical", "Direction", "Value"})

	for physical, logical := range pins {
		if p := this.GPIO.PhysicalPinForPin(logical); p > 0 {
			mode := this.GPIO.GetPinMode(logical)
			value := "-"
			if mode == gopi.GPIO_INPUT || mode == gopi.GPIO_OUTPUT {
				value = fmt.Sprint(this.GPIO.ReadPin(logical))
			}
			table.Append([]string{
				fmt.Sprintf("%v", p),
				fmt.Sprint(logical),
				fmt.Sprint(mode),
				value,
			})
		} else {
			table.Append([]string{
				fmt.Sprintf("%v", physical+1), "", "", "",
			})
		}
	}
	table.Render()

	// Return success
	return nil
}

func (this *app) GPIOParsePins(args []string) ([]gopi.GPIOPin, error) {
	var result []gopi.GPIOPin

	for _, arg := range args {
		if pin, err := strconv.ParseUint(arg, 0, 32); err == nil {
			if logical := this.GPIO.PhysicalPin(uint(pin)); logical == gopi.GPIO_PIN_NONE {
				return nil, gopi.ErrBadParameter.WithPrefix(arg)
			} else {
				result = append(result, logical)
			}
		} else if pin := rePin.FindStringSubmatch(arg); len(args) > 0 {
			pin_, _ := strconv.ParseUint(pin[2], 0, 32)
			logical := gopi.GPIOPin(pin_)
			if physical := this.GPIO.PhysicalPinForPin(logical); physical == 0 {
				return nil, gopi.ErrBadParameter.WithPrefix(arg)
			} else {
				result = append(result, logical)
			}
		} else {
			return nil, gopi.ErrBadParameter.WithPrefix(arg)
		}
	}

	// Return success
	return result, nil
}
