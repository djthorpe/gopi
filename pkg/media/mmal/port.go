// +build mmal

package mmal

import (
	"fmt"
	"io"

	"github.com/djthorpe/gopi/v3"
	mmal "github.com/djthorpe/gopi/v3/pkg/sys/mmal"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Port struct {
	port     *mmal.MMALPort
	pool     *mmal.MMALPool
	queue    *mmal.MMALQueue
	r        io.Reader
	w        io.Writer
	eor, eow bool
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewReaderPort(r io.Reader, ctx *mmal.MMALPort) (*Port, error) {
	if r == nil || ctx == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("NewReaderPort")
	} else if this, err := NewDataPort(ctx); err != nil {
		return nil, err
	} else {
		this.r = r
		return this, nil
	}
}

func NewWriterPort(w io.Writer, ctx *mmal.MMALPort) (*Port, error) {
	if w == nil || ctx == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("NewWriterPort")
	} else if this, err := NewDataPort(ctx); err != nil {
		return nil, err
	} else if queue := mmal.MMALQueueCreate(); queue == nil {
		this.port.FreePool(this.pool)
		return nil, gopi.ErrInternalAppError.WithPrefix("NewWriterPort")
	} else {
		this.w = w
		this.queue = queue
		return this, nil
	}
}

func NewDataPort(ctx *mmal.MMALPort) (*Port, error) {
	this := new(Port)

	if ctx == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("NewDataPort")
	} else {
		this.port = ctx
	}

	if pool := ctx.CreatePool(ctx.BufferMin()); pool == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewDataPort")
	} else {
		this.pool = pool
	}

	return this, nil
}

func (this *Port) Dispose() error {
	var result error

	// If port enabled, then disable it
	if this.port.Enabled() {
		if err := this.port.Disable(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Free resources
	this.port.FreePool(this.pool)
	if this.queue != nil {
		this.queue.Free()
	}

	// Release resources
	this.port = nil
	this.pool = nil
	this.queue = nil
	this.r = nil
	this.w = nil

	// Return success
	return result
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *Port) Enable() error {
	if this.r != nil {
		return this.port.EnableWithCallback(this.reader_callback)
	} else if this.w != nil {
		return this.port.EnableWithCallback(this.writer_callback)
	} else {
		return this.port.Enable()
	}
}

func (this *Port) Disable() error {
	return this.port.Disable()
}

func (this *Port) Read() (uint, error) {
	var n uint
	var err error

	if this.r == nil || this.eor {
		return 0, nil
	}

	// Obtain a free buffer from the pool and fill it
	buffer := this.pool.Get()
	if buffer == nil {
		return 0, nil
	} else if n, err = buffer.Fill(this.r); err == io.EOF {
		buffer.SetFlags(mmal.MMAL_BUFFER_HEADER_FLAG_EOS)
		this.eor = true
	} else if err != nil {
		return 0, err
	}

	// Send buffer to port
	if err := this.port.SendBuffer(buffer); err != nil {
		buffer.Release()
		return 0, err
	}

	// Return success
	return n, nil
}

func (this *Port) Write() (uint, error) {
	var n uint

	if this.w == nil || this.queue == nil {
		return 0, nil
	}

	// Send empty buffers to the output port of the decoder to allow the decoder to start
	// producing frames as soon as it gets input data
	for {
		if buffer := this.pool.Get(); buffer == nil {
			break
		} else if err := this.port.SendBuffer(buffer); err != nil {
			return n, err
		}
	}

	// Obtain free buffers from the queue to process
	for {
		buffer := this.queue.Get()
		if buffer == nil {
			return n, nil
		}
		this.Debug("Write: ", buffer)
		if buffer.HasFlags(mmal.MMAL_BUFFER_HEADER_FLAG_EOS) {
			this.Debug("  -> EOW")
			this.eow = true
		}
		if evt := buffer.Event(); evt == 0 {
			this.Debug("  -> WRITE ", buffer)
			if n_, err := this.w.Write(buffer.AsData()); err != nil {
				buffer.Release()
				return n, err
			} else {
				n = n + uint(n_)
			}
		} else if buffer.Event() == mmal.MMAL_EVENT_FORMAT_CHANGED {
			event := buffer.AsFormatChangeEvent()
			this.Debug("  FORMAT CHANGED ", event)
			/* RESIZE POOL HERE, COPY FORMAT */

			// Copy over the new format and re-enable the port
			if err := this.change_format(event); err != nil {
				buffer.Release()
				return n, err
			}
		} else {
			this.Debug("  UNHANDLED EVENT ", evt)
		}
	}

	// Return success
	return n, nil
}

func (this *Port) Debug(a ...interface{}) {
	fmt.Println(a...)
}

////////////////////////////////////////////////////////////////////////////////
// CALLBACKS

// reader_callback is called when a buffer should be discarded on an input port
func (this *Port) reader_callback(_ *mmal.MMALPort, buffer *mmal.MMALBuffer) {
	buffer.Release()
}

// writer_callback is called when a buffer should be placed in output queue
func (this *Port) writer_callback(_ *mmal.MMALPort, buffer *mmal.MMALBuffer) {
	this.queue.Put(buffer)
}

func (this *Port) change_format(event *mmal.MMALStreamFormatEvent) error {
	if err := this.port.Disable(); err != nil {
		return err
	} else if err := this.port.FormatFullCopy(event.Format()); err != nil {
		return err
	} else if err := this.port.FormatCommit(); err != nil {
		return err
	} else if err := this.port.Enable(); err != nil {
		return err
	}
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Port) Name() string {
	return this.port.Name()
}

func (this *Port) String() string {
	str := "<mmal.port"
	if this.r != nil {
		str += ".reader"
	}
	if this.w != nil {
		str += ".writer"
	}
	if this.eor {
		str += " EOR"
	}
	if this.eow {
		str += " EOW"
	}
	str += " port=" + fmt.Sprint(this.port)
	str += " pool=" + fmt.Sprint(this.pool)
	return str + ">"
}
