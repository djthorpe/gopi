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
	conn map[string]*Conn
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

	// Make map of devices and connections
	this.cast = make(map[string]*Cast)
	this.conn = make(map[string]*Conn)

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
	this.conn = nil

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

func (this *Manager) Connect(cast gopi.Cast) error {
	// Check for bad parameters
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("Connect")
	}

	// Device should have been discovered
	var result error
	if cast := this.getCastForId(cast.Id()); cast == nil {
		return gopi.ErrNotFound.WithPrefix("Connect")
	} else if conn := this.getConnForId(cast.id); conn != nil {
		return gopi.ErrOutOfOrder.WithPrefix("Connect")
	} else if conn, err := cast.Connect(serviceConnectTimeout); err != nil {
		return err
	} else {
		this.setConnForId(cast.id, conn)

		// Emit connect message
		if err := this.Publisher.Emit(NewCastEvent(cast, gopi.CAST_FLAG_CONNECT), false); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

func (this *Manager) Disconnect(cast gopi.Cast) error {
	// Check for bad parameters
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("Disconnect")
	}

	// Device should have been discovered
	var result error
	if cast := this.getCastForId(cast.Id()); cast == nil {
		return gopi.ErrNotFound.WithPrefix("Disconnect")
	} else if conn := this.getConnForId(cast.id); conn == nil {
		return gopi.ErrOutOfOrder.WithPrefix("Disconnect")
	} else {
		// Close connection
		this.setConnForId(cast.id, nil)
		if err := conn.Close(); err != nil {
			result = multierror.Append(result, err)
		}

		// Emit disconnect message
		if err := this.Publisher.Emit(NewCastEvent(cast, gopi.CAST_FLAG_DISCONNECT), false); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
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

	// Remove connection from list
	var result error
	if conn := this.getConnForId(cast.id); conn != nil {
		this.setConnForId(cast.id, nil)
		if err := conn.Close(); err != nil {
			result = multierror.Append(result, err)
		}
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

func (this *Manager) getConnForId(id string) *Conn {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if conn, exists := this.conn[id]; exists {
		return conn
	} else {
		return nil
	}
}

func (this *Manager) setConnForId(id string, conn *Conn) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if conn == nil {
		delete(this.conn, id)
	} else {
		this.conn[id] = conn
	}
}

// CastEvent returns any changes to a chromecast if it is already
// discovered or returns DISCOVERY flag
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
