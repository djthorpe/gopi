package googlecast

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	sync.RWMutex
	gopi.Unit
	gopi.ServiceDiscovery

	// Connected Cast Devices
	dev map[string]*Cast
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	serviceTypeCast       = "_googlecast._tcp"
	serviceConnectTimeout = time.Second * 15
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	if this.ServiceDiscovery == nil {
		return gopi.ErrInternalAppError.WithPrefix("ServiceDiscovery")
	}

	// Make map of devices
	this.dev = make(map[string]*Cast)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Disconnect devices
	var result error
	for _, cast := range this.dev {
		if err := this.disconnect(cast); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.dev = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Manager) Devices(ctx context.Context) ([]gopi.Cast, error) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Perform the lookup
	records, err := this.ServiceDiscovery.Lookup(ctx, serviceTypeCast)
	if err != nil {
		return nil, err
	}

	result := make([]gopi.Cast, 0, len(records))
	for _, record := range records {
		if cast := NewCastFromRecord(record); cast == nil {
			continue
		} else if connected, exists := this.dev[cast.id]; exists {
			result = append(result, connected)
		} else {
			result = append(result, cast)
		}
	}

	// Return success
	return result, nil
}

func (this *Manager) Connect(device gopi.Cast) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check for bad parameters
	if device == nil {
		return gopi.ErrBadParameter.WithPrefix("Connect")
	}

	// Check for already connected
	key := device.Id()
	if _, exists := this.dev[key]; exists {
		return gopi.ErrDuplicateEntry.WithPrefix("Connect")
	}

	// Do the connection
	if device_, ok := device.(*Cast); ok == false {
		return gopi.ErrInternalAppError.WithPrefix("Connect")
	} else if err := this.connect(device_); err != nil {
		return err
	} else {
		this.dev[key] = device_
	}

	// Return success
	return nil
}

func (this *Manager) Disconnect(device gopi.Cast) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error
	if device == nil {
		return gopi.ErrBadParameter.WithPrefix("Disconnect")
	}

	key := device.Id()
	if connected, exists := this.dev[key]; exists == false {
		return gopi.ErrNotFound.WithPrefix("Disconnect")
	} else if err := this.disconnect(connected); err != nil {
		result = multierror.Append(result, err)
	}
	delete(this.dev, key)
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<cast.manager"
	for _, device := range this.dev {
		str += fmt.Sprint(" %v=%v", device.Id(), device)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) disconnect(device *Cast) error {
	return device.Disconnect()
}

func (this *Manager) connect(device *Cast) error {
	return device.ConnectWithTimeout(serviceConnectTimeout)
}
