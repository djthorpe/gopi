package mdns

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/miekg/dns"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Discovery struct {
	gopi.Unit
	gopi.Publisher
	*Listener
	sync.WaitGroup
}

const (
	queryServices = "_services._dns-sd._udp"
	queryRepeat   = 2
	queryBackoff  = time.Millisecond * 50
)

///////////////////////////////////////////////////////////////////////////////
// RUN

func (this *Discovery) Run(ctx context.Context) error {
	// Subscribe to DNS messages
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

FOR_LOOP:
	for {
		select {
		case evt := <-ch:
			if msg, ok := evt.(*msgevent); ok {
				if services := NewServices(msg.Msg, this.Listener.Domain()).Services(); len(services) > 0 {
					for _, service := range services {
						if err := this.Publisher.Emit(service, true); err != nil {
							this.Print(err)
						}
					}
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

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Discovery) Lookup(ctx context.Context, srv string) ([]gopi.ServiceRecord, error) {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	// Collect services in goroutine
	var wg sync.WaitGroup
	ch := this.Publisher.Subscribe()
	records := make([]*service, 0, 10)

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
					if service.Service() == srv {
						records = append(records, service)
					}
				}
			}
		}
	}()

	// Query for lookup on all interfaces
	zone := this.Listener.Domain()
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
						key := srv.Name()
						names[key] = true
					}
				}
			}
		}
	}()

	// Query for services on all interfaces
	zone := this.Listener.Domain()
	if err := this.query(ctx, msgQueryServices(zone), 0); err != nil {
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

func (this *Discovery) Serve(context.Context, []gopi.ServiceRecord) error {
	return gopi.ErrNotImplemented
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

func msgQueryServices(zone string) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetQuestion(fqn(queryServices)+fqn(zone), dns.TypePTR)
	msg.RecursionDesired = false
	return msg
}

func msgQueryLookup(srv, zone string) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetQuestion(fqn(srv)+fqn(zone), dns.TypePTR)
	msg.RecursionDesired = false
	return msg
}

func fqn(value string) string {
	return strings.Trim(value, ".") + "."
}
