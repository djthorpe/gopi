// +build drm

package display

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Platform
	sync.RWMutex

	node    *string
	fh      *os.File
	res     *drm.ModeResources
	display map[uint32]*Display
}

// DisplayManager is the interface for DisplayManager with extra
// properties needed for SurfaceManager
type DisplayManager interface {
	gopi.DisplayManager

	Fd() uintptr
	Width() (uint32, uint32)
	Height() (uint32, uint32)
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) Define(cfg gopi.Config) error {
	this.node = cfg.FlagString("display.gpu", "", "GPU device")
	return nil
}

func (this *Manager) New(gopi.Config) error {
	if this.Platform == nil {
		return gopi.ErrInternalAppError.WithPrefix("Invalid gopi.Platform unit")
	}
	if *this.node == "" {
		*this.node = this.chooseGpu()
	}
	if *this.node == "" {
		return gopi.ErrBadParameter.WithPrefix("Missing flag -display.gpu")
	}
	if fh, err := drm.OpenDevice(*this.node); err != nil {
		return err
	} else {
		this.fh = fh
	}
	if res, err := drm.GetResources(this.fh.Fd()); err != nil {
		return err
	} else {
		this.res = res
	}

	// Create display map
	this.display = make(map[uint32]*Display)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Release display resources
	for _, display := range this.display {
		if display == nil {
			continue
		} else if err := display.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release DRM resources
	if this.res != nil {
		this.res.Free()
	}

	// Release device filehandle
	if this.fh != nil {
		if err := this.fh.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.fh = nil
	this.res = nil
	this.display = nil

	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Manager) Fd() uintptr {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.fh != nil {
		return this.fh.Fd()
	} else {
		return 0
	}
}

func (this *Manager) Width() (uint32, uint32) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.res != nil {
		return this.res.Width()
	} else {
		return 0, 0
	}

}

func (this *Manager) Height() (uint32, uint32) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.res != nil {
		return this.res.Height()
	} else {
		return 0, 0
	}

}

////////////////////////////////////////////////////////////////////////////////
// DISPLAYS

// GetDisplay returns display matched by Id
func (this *Manager) Display(id uint32) gopi.Display {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check parameters
	if this.fh == nil {
		return nil
	}

	// Return previously opened display
	display, exists := this.display[id]
	if exists {
		return display
	}

	conns := this.connectors(func(conn *drm.ModeConnector) bool {
		if id == 0 {
			return conn.Status() == drm.ModeConnectionConnected
		} else {
			return conn.Id() == id
		}
	})

	// Check condition where no connectors
	if len(conns) == 0 {
		return nil
	}

	// Free unused connectors
	for i, ctx := range conns {
		if i > 0 {
			ctx.Free()
		}
	}

	// Get encoder and CRTC
	encoder, err := drm.GetEncoder(this.fh.Fd(), conns[0].Encoder())
	if err != nil {
		conns[0].Free()
		return nil
	}
	crtc, err := drm.GetCRTC(this.fh.Fd(), encoder.Crtc())
	if err != nil {
		encoder.Free()
		conns[0].Free()
		return nil
	}

	// Make display
	if display = NewDisplay(conns[0], encoder, crtc); display == nil {
		crtc.Free()
		encoder.Free()
		conns[0].Free()
		return nil
	} else {
		key := conns[0].Id()
		this.display[key] = display
	}

	// Return display
	return display
}

func (this *Manager) PrimaryDisplay() gopi.Display {
	return this.Display(0)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<displaymanager.drm"
	if *this.node != "" {
		str += " gpu=" + strconv.Quote(*this.node)
	}
	if this.fh != nil {
		str += " fd=" + fmt.Sprint(this.fh.Fd())
	}
	if this.res != nil {
		str += " resources=" + fmt.Sprint(this.res)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) chooseGpu() string {
	if t := this.Platform.Type(); t&gopi.PLATFORM_RPI == 0 {
		// Find first node
		if nodes := drm.Devices(); len(nodes) > 0 {
			return nodes[0]
		} else {
			return ""
		}
	} else if t&gopi.PLATFORM_BCM2838_ARM8 != 0 {
		// Raspberry Pi 4
		return "card1"
	} else {
		// Other Raspberry Pi
		return "card0"
	}
}

func (this *Manager) connectors(filter func(*drm.ModeConnector) bool) []*drm.ModeConnector {
	connectors := this.res.Connectors()
	result := make([]*drm.ModeConnector, 0, len(connectors))
	for _, id := range connectors {
		if ctx, err := drm.GetConnector(this.fh.Fd(), id); err != nil {
			continue
		} else if filter != nil && filter(ctx) == false {
			ctx.Free()
		} else {
			result = append(result, ctx)
		}
	}
	return result
}
