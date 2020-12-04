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
	queryRepeat   = 3
	queryBackoff  = time.Millisecond * 200
)

///////////////////////////////////////////////////////////////////////////////
// RUN

func (this *Discovery) Run(ctx context.Context) error {
	// Subscribe to messages
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

FOR_LOOP:
	for {
		select {
		case evt := <-ch:
			if msg, ok := evt.(*dnsevent); ok {
				NewServices(msg.Msg(), this.Listener.Domain())
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

func (this *Discovery) EnumerateServices(ctx context.Context) error {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	// Query for services on all interfaces
	zone := this.Listener.Domain()
	if err := this.query(ctx, msgQueryServices(zone), 0); err != nil {
		return err
	}

	// Wait for completion
	<-ctx.Done()

	return ctx.Err()
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

func msgQueryServices(domain string) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetQuestion(fqn(queryServices)+fqn(domain), dns.TypePTR)
	msg.RecursionDesired = false
	return msg
}

func fqn(value string) string {
	return strings.Trim(value, ".") + "."
}
