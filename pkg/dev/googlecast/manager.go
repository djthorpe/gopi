package googlecast

import (
	"context"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.ServiceDiscovery
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	serviceTypeCast = "_googlecast._tcp"
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	if this.ServiceDiscovery == nil {
		return gopi.ErrInternalAppError.WithPrefix("ServiceDiscovery")
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Manager) Devices(ctx context.Context) ([]gopi.Cast, error) {
	// Perform the lookup
	records, err := this.ServiceDiscovery.Lookup(ctx, serviceTypeCast)
	if err != nil {
		return nil, err
	}

	result := make([]gopi.Cast, 0, len(records))
	for _, record := range records {
		if cast := NewCastFromRecord(record); cast != nil {
			result = append(result, cast)
		}
	}

	// Return success
	return result, nil
}

func (this *Manager) Connect(gopi.Cast) error {
	return gopi.ErrNotImplemented
}

func (this *Manager) Disconnect(gopi.Cast) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<cast.manager"
	return str + ">"
}
