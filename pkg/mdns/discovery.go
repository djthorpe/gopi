package mdns

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
	"github.com/miekg/dns"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Discovery struct {
	gopi.Unit
	gopi.Publisher
	sync.WaitGroup
	gopi.Logger
	*Listener
	*Responder
}

const (
	queryServices   = "_services._dns-sd._udp"
	queryRepeat     = 0
	queryBackoff    = time.Millisecond * 250
	queryDefaultTTL = 60 * 30 // In seconds (30 mins)
)

///////////////////////////////////////////////////////////////////////////////
// RUN

func (this *Discovery) Run(ctx context.Context) error {
	if this.Publisher == nil {
		return gopi.ErrInternalAppError.WithPrefix("Missing gopi.Publisher")
	}

	// Subscribe to DNS messages
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

FOR_LOOP:
	for {
		select {
		case evt := <-ch:
			if msg, ok := evt.(*msgevent); ok {
				if err := this.ParseEmit(msg.Msg); err != nil {
					this.Print(err)
				}
			}
		case <-ctx.Done():
			break FOR_LOOP
		}
	}

	// Wait for EnumererateServices to complete
	this.WaitGroup.Wait()

	// Return context state
	return ctx.Err()
}

func (this *Discovery) ParseEmit(msg *dns.Msg) error {
	// Parse into services
	services := NewServices(msg, this.Listener.Zone()).Services()
	if len(services) == 0 {
		return nil
	}

	// Emit services
	var result error
	for _, service := range services {
		if err := this.Publisher.Emit(service, true); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Discovery) Lookup(ctx context.Context, srv string) ([]gopi.ServiceRecord, error) {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	// Sanitize srv
	if srv == "" {
		return nil, gopi.ErrBadParameter.WithPrefix(srv)
	} else {
		srv = fqn(srv)
	}

	// Collect services in goroutine
	var wg sync.WaitGroup
	ch := this.Publisher.Subscribe()
	records := make(map[string]*service, 10)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer this.Publisher.Unsubscribe(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case evt := <-ch:
				if service, ok := evt.(*service); ok {
					if service.Service() == srv && service.ttl != 0 {
						key := service.Instance()
						records[key] = service
					}
				}
			}
		}
	}()

	// Query for lookup on all interfaces
	zone := this.Listener.Zone()
	if err := this.query(ctx, msgQueryLookup(srv, zone), 0); err != nil {
		return nil, err
	}

	// Wait for end of collection of names
	wg.Wait()

	// Collect services
	result := make([]gopi.ServiceRecord, 0, len(records))
	for _, record := range records {
		result = append(result, record)
	}

	// Return result
	return result, nil
}

func (this *Discovery) EnumerateServices(ctx context.Context) ([]string, error) {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	// Collect names in goroutine
	var wg sync.WaitGroup
	names := make(map[string]bool)
	query := msgQueryServices(this.Listener.Zone())
	ch := this.Publisher.Subscribe()
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer this.Publisher.Unsubscribe(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case evt := <-ch:
				if srv, ok := evt.(*service); ok {
					if srv.Service() == fqn(queryServices) && srv.ttl != 0 {
						key := fqn(srv.Name())
						names[key] = true
					}
				}
			}
		}
	}()

	// Query for services on all interfaces
	if err := this.query(ctx, query, 0); err != nil {
		return nil, err
	}

	// Wait for end of collection of names
	wg.Wait()

	// Collect names
	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}

	return result, nil
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Query sends a message
func (this *Discovery) query(ctx context.Context, msg *dns.Msg, iface int) error {
	timer := time.NewTimer(1 * time.Nanosecond)
	defer timer.Stop()
	c := 0
	for {
		c++
		select {
		case <-timer.C:
			if err := this.Listener.Send(msg, iface); err != nil {
				return err
			} else if c >= queryRepeat {
				return nil
			}
			timer.Reset(queryBackoff * time.Duration(c))
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Return a query message looking up all services
func msgQueryServices(zone string) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetQuestion(fqn(queryServices)+fqn(zone), dns.TypePTR)
	msg.RecursionDesired = false
	return msg
}

// Return a query message looking up a specific service record
func msgQueryLookup(srv, zone string) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetQuestion(fqn(srv)+fqn(zone), dns.TypePTR)
	msg.RecursionDesired = false
	return msg
}

// Return fully-qualified value
func fqn(value string) string {
	return strings.Trim(value, ".") + "."
}

// Transform from fully-qualified value
func unfqn(value string) string {
	return strings.TrimSuffix(value, ".")
}
