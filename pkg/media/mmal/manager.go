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

	c map[string]*mmal.MMALComponent
	p []*Port
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Logger)

	// Create component map
	this.c = make(map[string]*mmal.MMALComponent, 10)

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
		if err := c.Free(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.c = nil
	this.p = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Manager) VideoDecoder() (*mmal.MMALComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_DECODER)
}

func (this *Manager) VideoEncoder() (*mmal.MMALComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_ENCODER)
}

func (this *Manager) VideoRenderer() (*mmal.MMALComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER)
}

func (this *Manager) ImageDecoder() (*mmal.MMALComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_IMAGE_DECODER)
}

func (this *Manager) ImageEncoder() (*mmal.MMALComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_IMAGE_ENCODER)
}

func (this *Manager) Camera() (*mmal.MMALComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_CAMERA)

}

func (this *Manager) CameraInfo() (*mmal.MMALComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_CAMERA_INFO)

}

func (this *Manager) VideoSplitter() (*mmal.MMALComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_VIDEO_SPLITTER)

}

func (this *Manager) AudioRenderer() (*mmal.MMALComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_AUDIO_RENDERER)

}

func (this *Manager) Clock() (*mmal.MMALComponent, error) {
	return this.component(mmal.MMAL_COMPONENT_DEFAULT_CLOCK)
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

// CreateReaderForComponent creates a data input, which reads data from elsewhere
func (this *Manager) CreateReaderForComponent(r io.Reader, component *mmal.MMALComponent, index uint) (*Port, error) {
	if component == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateReaderForComponent")
	} else if ports := component.InputPorts(); index >= uint(len(ports)) {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateReaderForComponent")
	} else if port, err := NewReaderPort(r, ports[index]); err != nil {
		return nil, err
	} else {
		this.RWMutex.Lock()
		defer this.RWMutex.Unlock()

		this.p = append(this.p, port)
		return port, nil
	}
}

// CreateOutputForComponent creates a data output, which sends data out of MMAL
func (this *Manager) CreateWriterForComponent(w io.Writer, component *mmal.MMALComponent, index uint) (*Port, error) {
	if component == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateWriterForComponent")
	} else if ports := component.OutputPorts(); index >= uint(len(ports)) {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateWriterForComponent")
	} else if port, err := NewWriterPort(w, ports[index]); err != nil {
		return nil, err
	} else {
		this.RWMutex.Lock()
		defer this.RWMutex.Unlock()

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

	// Run loop
FOR_LOOP:
	for {
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

func (this *Manager) component(name string) (*mmal.MMALComponent, error) {
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
		this.c[name] = c
	}

	// Return component
	return this.c[name], nil
}

func (this *Manager) get(name string) *mmal.MMALComponent {
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
