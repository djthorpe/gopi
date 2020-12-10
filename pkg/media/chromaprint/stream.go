// +build chromaprint

package chromaprint

import (
	"fmt"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	chromaprint "github.com/djthorpe/gopi/v3/pkg/sys/chromaprint"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type stream struct {
	ctx         *chromaprint.Context
	fingerprint string
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewStream(rate, channels int) (*stream, error) {
	this := new(stream)

	if ctx := chromaprint.NewChromaprint(chromaprint.ALGORITHM_DEFAULT); ctx == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewStream")
	} else if err := ctx.Start(rate, channels); err != nil {
		ctx.Free()
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *stream) Close() error {
	this.ctx.Free()
	this.ctx = nil
	this.fingerprint = ""
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

func (this *stream) Write(data []int16) error {
	return this.ctx.Feed(data)
}

func (this *stream) GetFingerprint() (string, error) {
	if this.fingerprint != "" {
		return this.fingerprint, nil
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
	return this.ctx.DurationMs()
}

func (this *stream) Channels() int {
	return this.ctx.Channels()
}

func (this *stream) Rate() int {
	return this.ctx.Rate()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *stream) String() string {
	str := "<chromaprint.stream"
	str += " context=" + fmt.Sprint(this.ctx)
	return str + ">"
}
