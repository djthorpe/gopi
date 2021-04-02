package chromecast

import (
	"context"
	"fmt"
	"sync"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	sync.RWMutex
	gopi.Unit
	gopi.ServiceDiscovery
	gopi.Publisher
	gopi.Logger

	cast map[string]*Cast
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
	this.Require(this.ServiceDiscovery, this.Logger, this.Publisher)

	// Make map of devices
	this.cast = make(map[string]*Cast)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Disconnect devices
	var result error
	for _, cast := range this.cast {
		if err := this.disconnect(cast); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.cast = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<cast.manager"
	for _, cast := range this.cast {
		str += fmt.Sprint(" ", cast)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) Devices(ctx context.Context) ([]gopi.Cast, error) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Perform the lookup
	records, err := this.ServiceDiscovery.Lookup(ctx, serviceTypeCast)
	if err != nil {
		return nil, err
	}

	// Return any casts found
	result := make([]gopi.Cast, 0, len(records))
	for _, record := range records {
		if cast := NewCastFromRecord(record); cast == nil {
			continue
		} else if connected, exists := this.cast[cast.id]; exists {
			result = append(result, connected)
		} else {
			result = append(result, cast)
		}
	}

	// Return success
	return result, nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) disconnect(cast *Cast) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Remove cast from list
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("Disconnect")
	} else if _, exists := this.cast[cast.id]; exists == false {
		return gopi.ErrNotFound.WithPrefix("Disconnect")
	} else {
		delete(this.cast, cast.id)
	}

	// Call disconnect for the chromecast
	var result error
	if err := cast.Disconnect(); err != nil {
		result = multierror.Append(result, err)
	}

	// Return any errors
	return result
}
