package mdns

import (
	"context"
	"fmt"
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
	queryServices        = "_services._dns-sd._udp"
	queryRepeat          = 2
	queryBackoffDuration = time.Millisecond * 100
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
			fmt.Println("MSG", evt)
		case <-ctx.Done():
			break FOR_LOOP
		}
	}

	// Wait for EnumererateServices to complete
	this.WaitGroup.Wait()

	// Return success
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Discovery) EnumerateServices(ctx context.Context) error {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	// Query for services
	zone := this.Listener.Domain()
	return this.query(ctx, msgQueryServices(zone), 0)
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Query sends a message
func (this *Discovery) query(ctx context.Context, msg *dns.Msg, iface int) error {
	ticker := time.NewTimer(1 * time.Nanosecond)
	defer ticker.Stop()

	for i := 1; i <= queryRepeat; i++ {
		select {
		case <-ticker.C:
			if err := this.Listener.Send(msg, iface); err != nil {
				return err
			}
			ticker.Reset(time.Duration(i) * queryBackoffDuration)
		case <-ctx.Done():
			ticker.Stop()
			return ctx.Err()
		}
	}

	// Return success
	return nil
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
