// +build mmal

package mmal

import (
	"context"
	"io"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	mmal "github.com/djthorpe/gopi/v3/pkg/sys/mmal"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Logger
	sync.RWMutex

	c map[string]*component
	p []*Port

	// Channels for messages
	debug chan string
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Logger)

	// Create component map
	this.c = make(map[string]*component, 10)

	// Create debug channel
	this.debug = make(chan string)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Free ports
	for _, p := range this.p {
		if err := p.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Free componentns
	for _, c := range this.c {
		if err := c.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Close channels
	close(this.debug)

	// Release resources
	this.c = nil
	this.p = nil
	this.debug = nil

	// Return any errors
	return result
}

func (this *Manager) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-this.debug:
			this.Debug(msg)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Manager) VideoDecoder() (VideoComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_DECODER)
}

func (this *Manager) VideoEncoder() (VideoComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_ENCODER)
}

func (this *Manager) VideoRenderer() (VideoComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER)
}

func (this *Manager) ImageDecoder() (ImageComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_IMAGE_DECODER)
}

func (this *Manager) ImageEncoder() (ImageComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_IMAGE_ENCODER)
}

func (this *Manager) Camera() (CameraComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_CAMERA)
}

func (this *Manager) CameraInfo() (Component, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_CAMERA_INFO)
}

func (this *Manager) VideoSplitter() (VideoComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_SPLITTER)
}

func (this *Manager) AudioRenderer() (AudioComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_AUDIO_RENDERER)
}

func (this *Manager) Clock() (Component, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_CLOCK)
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

// CreateReaderForComponent creates a data input, which reads data from elsewhere
func (this *Manager) CreateReaderForComponent(r io.Reader, c Component, index uint) (*Port, error) {
	if c == nil || r == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateReaderForComponent")
	} else if ports := c.(*component).ctx.InputPorts(); index >= uint(len(ports)) {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateReaderForComponent")
	} else if port, err := NewReaderPort(r, ports[index]); err != nil {
		return nil, err
	} else {
		this.RWMutex.Lock()
		defer this.RWMutex.Unlock()

		// Set debug
		port.debug = this.debug

		this.p = append(this.p, port)
		return port, nil
	}
}

// CreateOutputForComponent creates a data output, which sends data out of MMAL
func (this *Manager) CreateWriterForComponent(w io.Writer, c Component, index uint) (*Port, error) {
	if c == nil || w == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateWriterForComponent")
	} else if ports := c.(*component).ctx.OutputPorts(); index >= uint(len(ports)) {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateWriterForComponent")
	} else if port, err := NewWriterPort(w, ports[index]); err != nil {
		return nil, err
	} else {
		this.RWMutex.Lock()
		defer this.RWMutex.Unlock()

		// Set debug
		port.debug = this.debug

		this.p = append(this.p, port)
		return port, nil
	}
}

func (this *Manager) Exec(ctx context.Context) error {
	// Enable ports
	for _, p := range this.p {
		if err := p.Enable(); err != nil {
			return err
		}
	}
	// Enable componentns
	for _, c := range this.c {
		if err := c.Enable(); err != nil {
			return err
		}
	}

	// Conditions for end of loop is eor=true and eow=true
	var eor, eow bool

	// Run loop - require both eor and eow to be true
FOR_LOOP:
	for eor == false || eow == false {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		default:
			// Read from inputs
			for _, p := range this.p {
				if n, err := p.Read(); err != nil {
					this.Print(p.Name(), " Read: ", err)
				} else if n > 0 {
					this.Print(p.Name(), " Read: ", n, " bytes")
				}
			}

			// Write to outputs
			for _, p := range this.p {
				if n, err := p.Write(); err != nil {
					this.Print(p.Name(), " Write: ", err)
				} else if n > 0 {
					this.Print(p.Name(), " Write: ", n, " bytes")
				}
			}

			// Check for EOR and EOW - need all EOR and EOW to be true
			// in order to finish the loop
			eor, eow = true, true
			for _, p := range this.p {
				if p.EOR() == false {
					eor = false
				}
				if p.EOW() == false {
					eow = false
				}
			}
			this.Print("EOR=", eor, " EOW=", eow)
		}
	}

	// Shutdown
	var result error

	// Append context error
	if err := ctx.Err(); err != nil {
		result = multierror.Append(result, err)
	}

	// Disable components
	for _, c := range this.c {
		if err := c.Disable(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Disable ports
	for _, p := range this.p {
		if err := p.Disable(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) component(name string) (*component, error) {
	// Attempt to get existing compnent
	if c := this.get(name); c != nil {
		return c, nil
	}

	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Create component
	if c, err := mmal.MMALComponentCreate(name); err != nil {
		return nil, err
	} else {
		this.c[name] = NewComponent(c)
	}

	// Return component
	return this.c[name], nil
}

func (this *Manager) get(name string) *component {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if c, exists := this.c[name]; exists {
		return c
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<mmal.manager"
	return str + ">"
}
