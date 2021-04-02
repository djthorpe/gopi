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
	serviceTypeCast       = "_googlecast._tcp."
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

func (this *Manager) Run(ctx context.Context) error {
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	// Loop handling messages until done
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case evt := <-ch:
			// Fire cast change events
			if record, ok := evt.(gopi.ServiceRecord); ok {
				if record.Service() == serviceTypeCast {
					if cast := NewCastFromRecord(record); cast != nil {
						if flags := this.castevent(cast); flags != gopi.CAST_FLAG_NONE {
							if err := this.Publisher.Emit(NewCastEvent(cast, flags), false); err != nil {
								this.Print("Chromecast:", err)
							}
						}
					}
				}
			}
		}
	}
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
	// Perform the lookup
	records, err := this.ServiceDiscovery.Lookup(ctx, serviceTypeCast)
	if err != nil {
		return nil, err
	}

	// Return any casts found
	result := make([]gopi.Cast, 0, len(records))
	for _, record := range records {
		cast := NewCastFromRecord(record)
		if cast == nil {
			continue
		}

		// Add cast, emit event
		if existing := this.getCastForId(cast.id); existing == nil {
			this.castevent(cast)
		}

		// Append cast onto results
		result = append(result, this.getCastForId(cast.id))
	}

	// Return success
	return result, nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) disconnect(cast *Cast) error {
	// Check parameters
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("Disconnect")
	}

	// Remove cast from list
	if existing := this.getCastForId(cast.id); existing == nil {
		return gopi.ErrNotFound.WithPrefix("Disconnect")
	} else {
		this.setCastForId(cast.id, nil)
	}

	// Call disconnect for the chromecast
	var result error
	if err := cast.Disconnect(); err != nil {
		result = multierror.Append(result, err)
	}

	// Return any errors
	return result
}

func (this *Manager) getCastForId(id string) *Cast {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if cast, exists := this.cast[id]; exists {
		return cast
	} else {
		return nil
	}
}

func (this *Manager) setCastForId(id string, cast *Cast) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if cast == nil {
		delete(this.cast, id)
	} else {
		this.cast[id] = cast
	}
}

func (this *Manager) castevent(cast *Cast) gopi.CastFlag {
	if other := this.getCastForId(cast.id); other == nil {
		this.setCastForId(cast.id, cast)
		return gopi.CAST_FLAG_DISCOVERY
	} else if flags := other.Equals(cast); flags == gopi.CAST_FLAG_NONE {
		return flags
	} else {
		other.updateFrom(cast)
		return flags
	}
}
