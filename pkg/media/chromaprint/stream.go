// +build chromaprint

package chromaprint

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	chromaprint "github.com/djthorpe/gopi/v3/pkg/sys/chromaprint"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type stream struct {
	sync.Mutex

	ctx            *chromaprint.Context
	rate, channels int
	duration       time.Duration
	fingerprint    string
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewStream(rate, channels int) (*stream, error) {
	this := new(stream)

	if ctx := chromaprint.NewChromaprint(chromaprint.ALGORITHM_DEFAULT); ctx == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewStream")
	} else {
		this.ctx = ctx
		this.rate = rate
		this.channels = channels
	}

	// Return success
	return this, nil
}

func (this *stream) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.ctx != nil {
		this.ctx.Free()
	}

	// Release resources
	this.ctx = nil
	this.fingerprint = ""

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

func (this *stream) Write(data []int16) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(data) == 0 || len(data)%this.channels != 0 {
		return gopi.ErrBadParameter.WithPrefix("Write")
	} else if this.ctx == nil {
		return gopi.ErrOutOfOrder.WithPrefix("Write")
	} else if this.duration == 0 {
		if err := this.ctx.Start(this.rate, this.channels); err != nil {
			return err
		}
	}

	// Write data and update duration field
	if err := this.ctx.Feed(data); err != nil {
		return err
	} else {
		samples_per_second := time.Duration(this.rate) * time.Duration(this.channels)
		this.duration += time.Second * time.Duration(len(data)) / samples_per_second
	}

	// Return success
	return nil
}

func (this *stream) GetFingerprint() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.fingerprint != "" {
		return this.fingerprint, nil
	} else if this.duration == 0 || this.ctx == nil {
		return "", gopi.ErrOutOfOrder.WithPrefix("GetFingerprint")
	}

	if err := this.ctx.Finish(); err != nil {
		return "", err
	} else if fp, err := this.ctx.GetFingerprint(); err != nil {
		return "", err
	} else {
		this.fingerprint = fp
	}

	// Return success
	return this.fingerprint, nil
}

func (this *stream) Duration() time.Duration {
	return this.duration
}

func (this *stream) Channels() int {
	return this.channels
}

func (this *stream) Rate() int {
	return this.rate
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *stream) String() string {
	str := "<chromaprint.stream"
	if this.ctx != nil {
		str += " context=" + fmt.Sprint(this.ctx)
	}
	if r := this.Rate(); r != 0 {
		str += " sample_rate=" + fmt.Sprint(r)
	}
	if c := this.Channels(); c != 0 {
		str += " channels=" + fmt.Sprint(c)
	}
	if d := this.Duration(); d != 0 {
		str += " duration=" + fmt.Sprint(d)
	}
	if this.fingerprint != "" {
		str += " fingerprint=" + strconv.Quote(this.fingerprint)
	}
	return str + ">"
}
