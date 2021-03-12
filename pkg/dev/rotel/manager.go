package rotel

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
	term "github.com/pkg/term"
)

// Ref: https://www.rotel.com/sites/default/files/product/rs232/A12-A14%20Protocol.pdf

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	sync.RWMutex
	gopi.Logger
	gopi.Publisher
	State

	// Flags
	tty  *string
	baud *uint

	fd  *term.Term // TTY file handle
	buf *strings.Builder
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEFAULT_TTY_BAUD    = 115200
	DEFAULT_TTY_TIMEOUT = 100 * time.Millisecond
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) Define(cfg gopi.Config) error {
	this.tty = cfg.FlagString("rotel.tty", "/dev/ttyUSB0", "RS232 device")
	this.baud = cfg.FlagUint("rotel.baud", DEFAULT_TTY_BAUD, "RS232 speed")
	return nil
}

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Publisher, this.Logger)

	// Check parameters
	if _, err := os.Stat(*this.tty); os.IsNotExist(err) {
		return gopi.ErrBadParameter.WithPrefix("-rotel.tty")
	} else if *this.baud == 0 {
		return gopi.ErrBadParameter.WithPrefix("-rotel.baud")
	}

	// Open term
	if fd, err := term.Open(*this.tty, term.Speed(int(*this.baud)), term.RawMode); err != nil {
		return err
	} else {
		this.fd = fd
		this.buf = new(strings.Builder)
	}

	// Set term read timeout
	if err := this.fd.SetReadTimeout(DEFAULT_TTY_TIMEOUT); err != nil {
		defer this.fd.Close()
		return err
	}

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Close RS232 connection
	if this.fd != nil {
		if err := this.fd.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Clear resources
	this.fd = nil
	this.buf = nil

	// Return any errors
	return result
}

func (this *Manager) Run(ctx context.Context) error {
	// Update rotel status every second
	timer := time.NewTimer(100 * time.Millisecond)
	defer timer.Stop()

	// Loop handling messages until done
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if cmd := this.State.Update(); cmd != "" {
				if err := this.writetty(cmd); err != nil {
					this.Print("writetty:  ", err)
				}
			}
			timer.Reset(time.Millisecond * 250)
		default:
			if err := this.readtty(); err != nil {
				this.Print("readtty: ", err)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) SetPower(state bool) error {
	if state {
		return this.writetty("power_on!")
	} else {
		return this.writetty("power_off!")
	}
}

func (this *Manager) SetSource(value string) error {
	switch value {
	case "pc_usb":
		return this.writetty("pcusb!")
	case "cd", "coax1", "coax2", "opt1", "opt2", "aux1", "aux2", "tuner", "photo", "usb", "bluetooth":
		return this.writetty(value + "!")
	default:
		return gopi.ErrBadParameter.WithPrefix("SetSource")
	}
}

func (this *Manager) SetVolume(value uint) error {
	if value < 1 || value > 96 {
		return gopi.ErrBadParameter.WithPrefix("SetVolume")
	} else {
		return this.writetty(fmt.Sprintf("vol_%02d!", value))
	}
}

func (this *Manager) SetMute(state bool) error {
	if state {
		return this.writetty("mute_on!")
	} else {
		return this.writetty("mute_off!")
	}
}

func (this *Manager) SetBypass(state bool) error {
	if state {
		return this.writetty("bypass_on!")
	} else {
		return this.writetty("bypass_off!")
	}
}

func (this *Manager) SetTreble(value int) error {
	if value < -10 || value > 10 {
		return gopi.ErrBadParameter.WithPrefix("SetTreble")
	} else if value == 0 {
		return this.writetty("treble_000!")
	} else if value < 0 {
		return this.writetty(fmt.Sprint("treble_", value, "!"))
	} else {
		return this.writetty(fmt.Sprint("treble_+", value, "!"))
	}
}

func (this *Manager) SetBass(value int) error {
	if value < -10 || value > 10 {
		return gopi.ErrBadParameter.WithPrefix("SetBass")
	} else if value == 0 {
		return this.writetty("bass_000!")
	} else if value < 0 {
		return this.writetty(fmt.Sprint("bass_", value, "!"))
	} else {
		return this.writetty(fmt.Sprint("bass_+", value, "!"))
	}
}

func (this *Manager) SetBalance(loc string, value uint) error {
	if value > 15 {
		return gopi.ErrBadParameter.WithPrefix("SetBalance")
	}
	switch loc {
	case "L", "R":
		if value == 0 {
			return this.writetty("balance_000!")
		} else {
			return this.writetty(fmt.Sprintf("balance_%v%02d!", loc, value))
		}
	default:
		return gopi.ErrBadParameter.WithPrefix("SetBalance")
	}
}

func (this *Manager) SetDimmer(value uint) error {
	if value > 6 {
		return gopi.ErrBadParameter.WithPrefix("SetDimmer")
	} else {
		return this.writetty(fmt.Sprint("dimmer_", value, "!"))
	}
}

func (this *Manager) Play() error {
	return this.writetty("play!")
}

func (this *Manager) Stop() error {
	return this.writetty("stop!")
}

func (this *Manager) Pause() error {
	return this.writetty("pause!")
}

func (this *Manager) NextTrack() error {
	return this.writetty("trkf!")
}

func (this *Manager) PrevTrack() error {
	return this.writetty("trkb!")
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<rotelmanager"
	str += fmt.Sprintf(" tty=%q", *this.tty)
	str += fmt.Sprint(" baud=", *this.baud)
	str += fmt.Sprint(" ", this.State.String())
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHDOS

func (this *Manager) readtty() error {
	n, err := this.fd.Available()
	if err != nil {
		return err
	} else if n == 0 {
		return nil
	}

	// Append data to the buffer and parse any parameters
	var result error
	var flags RotelFlag
	buf := make([]byte, n)
	if _, err := this.fd.Read(buf); err == io.EOF {
		return nil
	} else if err != nil {
		return err
	} else if _, err := this.buf.Write(buf); err != nil {
		return err
	} else if fields := strings.Split(this.buf.String(), "$"); len(fields) > 0 {
		// Parse each field and update state
		for _, param := range fields[0 : len(fields)-1] {
			if flag, err := this.State.Set(param); err != nil {
				result = multierror.Append(result, fmt.Errorf("%q: %w", param, err))
			} else {
				flags |= flag
			}
		}
		// Reset buffer with any remaining data not parsed
		this.buf.Reset()
		this.buf.WriteString(fields[len(fields)-1])
	}

	// If any flags set, then emit an event
	if flags != FLAG_NONE {
		if err := this.Emit(NewEvent(&this.State, flags), false); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

func (this *Manager) writetty(cmd string) error {
	this.Debugf("writetty: %q", cmd)
	_, err := this.fd.Write([]byte(cmd))
	return err
}
