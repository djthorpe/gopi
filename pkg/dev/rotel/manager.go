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
	"github.com/hashicorp/go-multierror"
	term "github.com/pkg/term"
)

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
