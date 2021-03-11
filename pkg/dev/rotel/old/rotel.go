/*
	Rotel RS232 Control
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package rotel

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	event "github.com/djthorpe/gopi/util/event"
	rotel "github.com/djthorpe/rotel"
	term "github.com/pkg/term"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Rotel struct {
	TTY      string
	BaudRate uint
}

type driver struct {
	log gopi.Logger
	tty string
	fd  *term.Term
	buf string

	model string
	state rotel.RotelState

	event.Tasks
	event.Publisher
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

const (
	BAUD_RATE_DEFAULT = 115200
	READ_TIMEOUT      = 100 * time.Millisecond
)

var (
	reModel   = regexp.MustCompile("^model=(\\w+)$")
	rePower   = regexp.MustCompile("^power=(on|standby)$")
	reVolume  = regexp.MustCompile("^volume=(\\d+)$")
	reBass    = regexp.MustCompile("^bass=([\\+\\-]?\\d+)$")
	reTreble  = regexp.MustCompile("^treble=([\\+\\-]?\\d+)$")
	reBalance = regexp.MustCompile("^balance=([LR]?)(\\d+)$")
	reMute    = regexp.MustCompile("^mute=(on|off)$")
	reSource  = regexp.MustCompile("^source=(\\w+)$")
	reFreq    = regexp.MustCompile("^freq=(.+)$")
	reBypass  = regexp.MustCompile("^bypass=(on|off)$")
	reSpeaker = regexp.MustCompile("^speaker=(a|b|a_b|off)$")
	reDimmer  = regexp.MustCompile("^dimmer=(\\d+)$")
	reUpdate  = regexp.MustCompile("^update_mode=(auto|manual)$")
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Rotel) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("<rotel.Open>{ config=%+v }", config)

	this := new(driver)
	this.log = logger

	// Check TTY parameter
	if config.TTY == "" {
		return nil, gopi.ErrBadParameter
	} else if _, err := os.Stat(config.TTY); os.IsNotExist(err) {
		return nil, fmt.Errorf("Open: %w", err)
	} else {
		this.tty = config.TTY
	}

	// Open term
	if fd, err := term.Open(this.tty, term.Speed(int(config.BaudRate)), term.RawMode); err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	} else {
		this.fd = fd
	}

	// Set term read timeout
	if err := this.fd.SetReadTimeout(READ_TIMEOUT); err != nil {
		defer this.fd.Close()
		return nil, fmt.Errorf("Open: %w", err)
	}

	// Start background thread
	this.Tasks.Start(this.ReadTask)

	// Success
	return this, nil
}

func (this *driver) Close() error {
	this.log.Debug("<rotel.Close>{ tty=%v state=%v }", strconv.Quote(this.tty), this.state)
	this.Lock()
	defer this.Unlock()

	// Remove subscribers
	this.Publisher.Close()

	// Stop background tasks
	if err := this.Tasks.Close(); err != nil {
		return err
	}

	// Close RS232 connection
	if this.fd != nil {
		if err := this.fd.Close(); err != nil {
			return err
		} else {
			this.fd = nil
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *driver) String() string {
	if this.fd == nil {
		return "<rotel>{ nil }"
	} else {
		return fmt.Sprintf("<rotel>{ tty=%v model=%v state=%v }", strconv.Quote(this.tty), strconv.Quote(this.model), this.state)
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET PARAMETERS

func (this *driver) Model() string {
	return this.model
}

func (this *driver) Get() rotel.RotelState {
	return this.state
}

func (this *driver) Set(state rotel.RotelState) error {
	this.log.Debug2("<rotel.Set>{ %v }", state)
	this.Lock()
	defer this.Unlock()

	if state.Power != rotel.ROTEL_POWER_NONE && state.Power != this.state.Power {
		if err := this.setPower(state.Power); err != nil {
			return fmt.Errorf("setPower: %w", err)
		}
		if this.state.Power == rotel.ROTEL_POWER_ON {
			// Power switching off so ignore other parameters
			return nil
		}
	}
	if state.Volume != rotel.ROTEL_VOLUME_NONE && state.Volume != this.state.Volume {
		this.log.Debug("current vol=%v set new vol=%v", this.state.Volume, state.Volume)
		if err := this.setVolume(state.Volume); err != nil {
			return fmt.Errorf("setVolume: %w", err)
		}
	}
	if state.Source != rotel.ROTEL_SOURCE_NONE && state.Source != this.state.Source {
		if err := this.setSource(state.Source); err != nil {
			return fmt.Errorf("setSource: %w", err)
		}
	}
	if state.Mute != rotel.ROTEL_MUTE_NONE && state.Mute != this.state.Mute {
		if err := this.setMute(state.Mute); err != nil {
			return fmt.Errorf("setMute: %w", err)
		}
	}
	if state.Bypass != rotel.ROTEL_BYPASS_NONE && state.Bypass != this.state.Bypass {
		if err := this.setBypass(state.Bypass); err != nil {
			return fmt.Errorf("setBypass: %w", err)
		}
	}
	if state.Treble != rotel.ROTEL_TONE_NONE && state.Treble != this.state.Treble {
		if err := this.setTreble(state.Treble); err != nil {
			return fmt.Errorf("setTreble: %w", err)
		}
	}
	if state.Bass != rotel.ROTEL_TONE_NONE && state.Bass != this.state.Bass {
		if err := this.setBass(state.Bass); err != nil {
			return fmt.Errorf("setBass: %w", err)
		}
	}
	if state.Balance != rotel.ROTEL_BALANCE_NONE && state.Balance != this.state.Balance {
		if err := this.setBalance(state.Balance); err != nil {
			return fmt.Errorf("setBalance: %w", err)
		}
	}
	if state.Dimmer != rotel.ROTEL_DIMMER_NONE && state.Dimmer != this.state.Dimmer {
		if err := this.setDimmer(state.Dimmer); err != nil {
			return fmt.Errorf("setDimmer: %w", err)
		}
	}
	if state.Update != rotel.ROTEL_UPDATE_NONE && state.Update != this.state.Update {
		if err := this.setUpdate(state.Update); err != nil {
			return fmt.Errorf("setUpdate: %w", err)
		}
	}
	if state.Speaker != rotel.ROTEL_SPEAKER_NONE && state.Speaker != this.state.Speaker {
		if err := this.setSpeaker(state.Speaker); err != nil {
			return fmt.Errorf("setSpeaker: %w", err)
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SET PARAMETERS

func (this *driver) setPower(value rotel.Power) error {
	this.log.Debug2("<rotel.SetPower>{ %v }", value)

	switch value {
	case rotel.ROTEL_POWER_ON:
		return this.write("power_on")
	case rotel.ROTEL_POWER_STANDBY:
		return this.write("power_off")
	default:
		return gopi.ErrBadParameter
	}
}

func (this *driver) setVolume(value rotel.Volume) error {
	this.log.Debug2("<rotel.SetVolume>{ %v }", value)

	if value >= rotel.ROTEL_VOLUME_MIN && value <= rotel.ROTEL_VOLUME_MAX {
		return this.write(fmt.Sprintf("vol_%d", value))
	} else {
		return gopi.ErrBadParameter
	}
}

func (this *driver) setSource(value rotel.Source) error {
	this.log.Debug2("<rotel.setSource>{ %v }", value)

	if str := sourceToString(value); str != "pc_usb" && str != "" {
		return this.write(str)
	} else if str == "pc_usb" {
		return this.write("pcusb")
	} else {
		return gopi.ErrBadParameter
	}
}

func (this *driver) setMute(value rotel.Mute) error {
	this.log.Debug2("<rotel.setMute>{ %v }", value)

	switch value {
	case rotel.ROTEL_MUTE_ON:
		return this.write("mute_on")
	case rotel.ROTEL_MUTE_OFF:
		return this.write("mute_off")
	default:
		return gopi.ErrBadParameter
	}
}

func (this *driver) setBypass(value rotel.Bypass) error {
	this.log.Debug2("<rotel.setBypass>{ %v }", value)

	switch value {
	case rotel.ROTEL_BYPASS_ON:
		return this.write("bypass_on")
	case rotel.ROTEL_BYPASS_OFF:
		return this.write("bypass_off")
	default:
		return gopi.ErrBadParameter
	}
}

func (this *driver) setSpeaker(value rotel.Speaker) error {
	this.log.Debug2("<rotel.setSpeaker>{ %v }", value)

	switch value {
	case rotel.ROTEL_SPEAKER_OFF:
		if err := this.write("speaker_a_off"); err != nil {
			return err
		} else if err := this.write("speaker_b_off"); err != nil {
			return err
		} else {
			return nil
		}
	case rotel.ROTEL_SPEAKER_A:
		if err := this.write("speaker_a_on"); err != nil {
			return err
		} else if err := this.write("speaker_b_off"); err != nil {
			return err
		} else {
			return nil
		}
	case rotel.ROTEL_SPEAKER_B:
		if err := this.write("speaker_b_on"); err != nil {
			return err
		} else if err := this.write("speaker_a_off"); err != nil {
			return err
		} else {
			return nil
		}
	case rotel.ROTEL_SPEAKER_ALL:
		if err := this.write("speaker_b_on"); err != nil {
			return err
		} else if err := this.write("speaker_a_on"); err != nil {
			return err
		} else {
			return nil
		}
	default:
		return gopi.ErrBadParameter
	}
}

func (this *driver) setBass(value rotel.Tone) error {
	this.log.Debug2("<rotel.setBass>{ %v }", value)

	if value >= rotel.ROTEL_TONE_MIN && value <= rotel.ROTEL_TONE_MAX {
		return this.write(fmt.Sprintf("bass_%d", value))
	} else if value == rotel.ROTEL_TONE_OFF {
		return this.write("bass_000")
	} else {
		return gopi.ErrBadParameter
	}
}

func (this *driver) setTreble(value rotel.Tone) error {
	this.log.Debug2("<rotel.setTreble>{ %v }", value)

	if value >= rotel.ROTEL_TONE_MIN && value <= rotel.ROTEL_TONE_MAX {
		return this.write(fmt.Sprintf("treble_%d", value))
	} else if value == rotel.ROTEL_TONE_OFF {
		return this.write("treble_000")
	} else {
		return gopi.ErrBadParameter
	}
}

func (this *driver) setBalance(value rotel.Balance) error {
	this.log.Debug2("<rotel.setBalance>{ %v }", value)

	if value < rotel.ROTEL_BALANCE_LEFT_MAX || value > rotel.ROTEL_BALANCE_RIGHT_MAX {
		return gopi.ErrBadParameter
	} else if value == rotel.ROTEL_BALANCE_OFF {
		return this.write("balance_000")
	} else if value >= rotel.ROTEL_BALANCE_LEFT_MAX {
		return this.write(fmt.Sprintf("balance_L%d", -value))
	} else if value <= rotel.ROTEL_BALANCE_RIGHT_MAX {
		return this.write(fmt.Sprintf("balance_R%d", value))
	} else {
		return gopi.ErrBadParameter
	}
}

func (this *driver) setDimmer(value rotel.Dimmer) error {
	this.log.Debug2("<rotel.setDimmer>{ %v }", value)

	if value == rotel.ROTEL_DIMMER_OFF {
		return this.write("dimmer_0")
	} else if value >= rotel.ROTEL_DIMMER_MIN && value <= rotel.ROTEL_DIMMER_MAX {
		return this.write(fmt.Sprintf("dimmer_%d", value))
	} else {
		return gopi.ErrBadParameter
	}
}

func (this *driver) setUpdate(value rotel.Update) error {
	this.log.Debug2("<rotel.setUpdate>{ %v }", value)

	if value == rotel.ROTEL_UPDATE_MANUAL {
		return this.write("rs232_update_off")
	} else if value == rotel.ROTEL_UPDATE_AUTO {
		return this.write("rs232_update_on")
	} else {
		return gopi.ErrBadParameter
	}
}

////////////////////////////////////////////////////////////////////////////////
// SEND COMMAND

func (this *driver) Send(value rotel.Command) error {
	this.log.Debug2("<rotel.Send>{ %v }", value)

	str := strings.TrimPrefix(fmt.Sprint(value), "ROTEL_COMMAND_")

	switch value {
	case rotel.ROTEL_COMMAND_PLAY, rotel.ROTEL_COMMAND_STOP, rotel.ROTEL_COMMAND_PAUSE:
		return this.write(strings.ToLower(str))
	case rotel.ROTEL_COMMAND_TRACK_NEXT:
		return this.write("trkf")
	case rotel.ROTEL_COMMAND_TRACK_PREV:
		return this.write("trkb")
	case rotel.ROTEL_COMMAND_MUTE_TOGGLE:
		return this.write("mute")
	case rotel.ROTEL_COMMAND_VOL_UP, rotel.ROTEL_COMMAND_VOL_DOWN:
		return this.write(strings.ToLower(str))
	case rotel.ROTEL_COMMAND_BASS_UP, rotel.ROTEL_COMMAND_TREBLE_UP, rotel.ROTEL_COMMAND_BASS_DOWN, rotel.ROTEL_COMMAND_TREBLE_DOWN:
		return this.write(strings.ToLower(str))
	case rotel.ROTEL_COMMAND_BASS_RESET:
		return this.write("bass_000")
	case rotel.ROTEL_COMMAND_TREBLE_RESET:
		return this.write("treble_000")
	case rotel.ROTEL_COMMAND_BALANCE_LEFT:
		return this.write("balance_l")
	case rotel.ROTEL_COMMAND_BALANCE_RIGHT:
		return this.write("balance_r")
	case rotel.ROTEL_COMMAND_BALANCE_RESET:
		return this.write("balance_000")
	case rotel.ROTEL_COMMAND_SPEAKER_A_TOGGLE:
		return this.write("speaker_a")
	case rotel.ROTEL_COMMAND_SPEAKER_B_TOGGLE:
		return this.write("speaker_b")
	case rotel.ROTEL_COMMAND_DIMMER_TOGGLE:
		return this.write("dimmer")
	case rotel.ROTEL_COMMAND_POWER_TOGGLE:
		return this.write("power_toggle")
	default:
		return gopi.ErrBadParameter
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *driver) write(command string) error {
	this.log.Debug2("<rotel.Write>{ %v! }", command)
	_, err := this.fd.Write([]byte(command + "!"))
	return err
}

func (this *driver) read(command string) error {
	this.log.Debug2("<rotel.Read>{ %v? }", command)
	_, err := this.fd.Write([]byte(command + "?"))
	return err
}

func (this *driver) parse(commands []string) error {
	this.log.Debug2("<rotel.Parse>{ %v }", strconv.Quote(strings.Join(commands, ",")))
	for _, command := range commands {
		if value := reModel.FindStringSubmatch(command); len(value) > 1 {
			this.model = value[1]
		} else if value := rePower.FindStringSubmatch(command); len(value) > 1 {
			switch value[1] {
			case "on":
				this.evtPower(rotel.ROTEL_POWER_ON)
			case "standby":
				this.evtPower(rotel.ROTEL_POWER_STANDBY)
			default:
				this.evtPower(rotel.ROTEL_POWER_NONE)
			}
		} else if value := reSource.FindStringSubmatch(command); len(value) > 1 {
			if source := stringToSource(value[1]); source != rotel.ROTEL_SOURCE_NONE {
				this.evtSource(source)
			} else {
				this.evtSource(rotel.ROTEL_SOURCE_OTHER)
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			}
		} else if value := reVolume.FindStringSubmatch(command); len(value) > 1 {
			if v, err := strconv.ParseUint(value[1], 10, 32); err == nil {
				if v >= uint64(rotel.ROTEL_VOLUME_MIN) && v <= uint64(rotel.ROTEL_VOLUME_MAX) {
					this.evtVolume(rotel.Volume(v))
				}
			} else {
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			}
		} else if value := reFreq.FindStringSubmatch(command); len(value) > 1 {
			if value[1] == "" {
				this.evtFreq("unknown")
			} else {
				this.evtFreq(value[1])
			}
		} else if value := reMute.FindStringSubmatch(command); len(value) > 1 {
			switch value[1] {
			case "on":
				this.evtMute(rotel.ROTEL_MUTE_ON)
			case "off":
				this.evtMute(rotel.ROTEL_MUTE_OFF)
			default:
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			}
		} else if value := reBypass.FindStringSubmatch(command); len(value) > 1 {
			switch value[1] {
			case "on":
				this.evtBypass(rotel.ROTEL_BYPASS_ON)
			case "off":
				this.evtBypass(rotel.ROTEL_BYPASS_OFF)
			default:
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			}
		} else if value := reBass.FindStringSubmatch(command); len(value) > 1 {
			if v, err := strconv.ParseInt(value[1], 10, 32); err == nil {
				if v == 0 {
					this.evtBass(rotel.ROTEL_TONE_OFF)
				} else if v >= int64(rotel.ROTEL_TONE_MIN) && v <= int64(rotel.ROTEL_TONE_MAX) {
					this.evtBass(rotel.Tone(v))
				}
			} else {
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			}
		} else if value := reTreble.FindStringSubmatch(command); len(value) > 1 {
			if v, err := strconv.ParseInt(value[1], 10, 32); err == nil {
				if v == 0 {
					this.evtTreble(rotel.ROTEL_TONE_OFF)
				} else if v >= int64(rotel.ROTEL_TONE_MIN) && v <= int64(rotel.ROTEL_TONE_MAX) {
					this.evtTreble(rotel.Tone(v))
				}
			} else {
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			}
		} else if value := reSpeaker.FindStringSubmatch(command); len(value) > 1 {
			switch value[1] {
			case "off":
				this.evtSpeaker(rotel.ROTEL_SPEAKER_OFF)
			case "a":
				this.evtSpeaker(rotel.ROTEL_SPEAKER_A)
			case "b":
				this.evtSpeaker(rotel.ROTEL_SPEAKER_B)
			case "a_b":
				this.evtSpeaker(rotel.ROTEL_SPEAKER_ALL)
			default:
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			}
		} else if value := reBalance.FindStringSubmatch(command); len(value) > 2 {
			if v, err := strconv.ParseUint(value[2], 10, 32); err != nil {
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			} else if value[1] == "" && v == 0 {
				this.evtBalance(rotel.ROTEL_BALANCE_OFF)
			} else if value[1] == "L" && v > 0 && v <= uint64(-rotel.ROTEL_BALANCE_LEFT_MAX) {
				this.evtBalance(rotel.Balance(-v))
			} else if value[1] == "R" && v > 0 && v <= uint64(rotel.ROTEL_BALANCE_RIGHT_MAX) {
				this.evtBalance(rotel.Balance(v))
			} else {
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			}
		} else if value := reDimmer.FindStringSubmatch(command); len(value) > 1 {
			if v_, err := strconv.ParseUint(value[1], 10, 32); err == nil {
				v := rotel.Dimmer(v_)
				if v == rotel.ROTEL_DIMMER_NONE {
					this.evtDimmer(rotel.ROTEL_DIMMER_OFF)
				} else if v >= rotel.ROTEL_DIMMER_MIN && v <= rotel.ROTEL_DIMMER_MAX {
					this.evtDimmer(v)
				}
			} else {
				return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
			}
		} else if value := reUpdate.FindStringSubmatch(command); len(value) > 1 {
			switch value[1] {
			case "auto":
				this.evtUpdate(rotel.ROTEL_UPDATE_AUTO)
			case "manual":
				this.evtUpdate(rotel.ROTEL_UPDATE_MANUAL)
			default:
				this.evtUpdate(rotel.ROTEL_UPDATE_OTHER)
			}
		} else {
			return fmt.Errorf("Cannot parse: %v", strconv.Quote(command))
		}
	}

	// Success
	return nil
}

func (this *driver) retrieveparams() error {
	if this.fd == nil {
		// If no file descriptor, do nothing
		return nil
	}

	switch {
	case this.model == "":
		return this.read("model")
	case this.state.Power == rotel.ROTEL_POWER_NONE:
		return this.read("power")
	case this.state.Power != rotel.ROTEL_POWER_ON:
		return nil
	case this.state.Update == rotel.ROTEL_UPDATE_NONE:
		if err := this.write("rs232_update_on"); err != nil {
			return err
		} else {
			this.evtUpdate(rotel.ROTEL_UPDATE_AUTO)
		}
	case this.state.Volume == rotel.ROTEL_VOLUME_NONE:
		return this.read("volume")
	case this.state.Source == rotel.ROTEL_SOURCE_NONE:
		return this.read("source")
	case this.state.Freq == "":
		return this.read("freq")
	case this.state.Bypass == rotel.ROTEL_BYPASS_NONE:
		return this.read("bypass")
	case this.state.Speaker == rotel.ROTEL_SPEAKER_NONE:
		return this.read("speaker")
	case this.state.Mute == rotel.ROTEL_MUTE_NONE:
		return this.read("mute")
	case this.state.Bass == rotel.ROTEL_TONE_NONE:
		return this.read("bass")
	case this.state.Treble == rotel.ROTEL_TONE_NONE:
		return this.read("treble")
	case this.state.Balance == rotel.ROTEL_BALANCE_NONE:
		return this.read("balance")
	case this.state.Dimmer == rotel.ROTEL_DIMMER_NONE:
		return this.read("dimmer")
	}

	// Nothing to do
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND TASKS

func (this *driver) ReadTask(start chan<- event.Signal, stop <-chan event.Signal) error {
	start <- gopi.DONE
	timer := time.NewTicker(5 * time.Second)

FOR_LOOP:
	for {
		select {
		case <-stop:
			break FOR_LOOP
		case <-timer.C:
			if err := this.retrieveparams(); err != nil {
				this.log.Warn("ReadTask: %v", err)
			}
		default:
			for {
				if n, err := this.fd.Available(); err != nil {
					this.log.Error("ReadTask: %v", err)
					break FOR_LOOP
				} else if n > 0 {
					buf := make([]byte, n)
					if _, err := this.fd.Read(buf); err != nil {
						this.log.Warn("ReadTask: %v", err)
					} else {
						// Append buffer and extract first field
						this.buf += string(buf)
						if fields := strings.Split(this.buf, "$"); len(fields) > 0 {
							this.buf = fields[len(fields)-1]
							if err := this.parse(fields[0 : len(fields)-1]); err != nil {
								this.log.Warn("ReadTask: %v", err)
							}
						}
					}
				} else {
					break
				}
			}
		}
	}

	// Stop timer
	timer.Stop()

	// Success
	return nil
}
