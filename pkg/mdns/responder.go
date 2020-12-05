package mdns

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	dns "github.com/miekg/dns"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Responder struct {
	gopi.Unit
	gopi.Publisher
	gopi.Logger
	sync.WaitGroup
	sync.RWMutex
	*Listener

	// Records to respond to
	names   []string
	records map[string][]gopi.ServiceRecord
}

// FuncServices returns fully-qualified service names and TTL (or zero for default)
type FuncServices func() ([]string, uint32)

// FuncRecordsForService returns service records for named service and TTL
type FuncRecordsForService func(string) ([]gopi.ServiceRecord, uint32)

///////////////////////////////////////////////////////////////////////////////
// RUN

func (this *Responder) Run(ctx context.Context) error {
	// Subscribe to DNS messages
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

FOR_LOOP:
	for {
		select {
		case evt := <-ch:
			if s, _ := this.Services(); len(s) == 0 {
				// Do not process messages where no services are defined
			} else if msg, ok := evt.(*msgevent); ok {
				if err := this.ProcessQuestion(msg); err != nil {
					this.Print(err)
				}
			}
		case <-ctx.Done():
			break FOR_LOOP
		}
	}

	// Wait for Serve to complete
	this.WaitGroup.Wait()

	// Return context state
	return ctx.Err()
}

///////////////////////////////////////////////////////////////////////////////
// GET PROPERTIES

func (this *Responder) Services() ([]string, uint32) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Return the current service names with unset TTL
	return this.names, 0
}

func (this *Responder) Records(name string) ([]gopi.ServiceRecord, uint32) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	key := fqn(name)
	if r, exists := this.records[key]; exists {
		return r, queryDefaultTTL
	} else {
		return nil, 0
	}
}

///////////////////////////////////////////////////////////////////////////////
// SET PROPERTIES

func (this *Responder) SetServices(r []gopi.ServiceRecord) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Map of services
	this.names = make([]string, 0, len(r))
	this.records = make(map[string][]gopi.ServiceRecord, len(r))
	for _, record := range r {
		key := fqn(record.Service())
		// Validate key
		if key == "." {
			continue
		}
		// Deal with new key
		if _, exists := this.records[key]; exists == false {
			this.records[key] = []gopi.ServiceRecord{}
			this.names = append(this.names, key)
		}
		// Append record to existing set
		this.records[key] = append(this.records[key], record)
	}

	// Return success
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Responder) ProcessQuestion(msg *msgevent) error {
	if msg.Msg == nil {
		return gopi.ErrBadParameter.WithPrefix("ProcessQuestion")
	}
	if msg.Opcode != dns.OpcodeQuery {
		return fmt.Errorf("Received query with non-zero Opcode (%v)", msg.Opcode)
	}
	if msg.Rcode != 0 {
		return fmt.Errorf("Received query with non-zero Rcode (%v)", msg.Rcode)
	}
	if msg.Truncated {
		return fmt.Errorf("DNS requests with high truncated bit not implemented")
	}

	// Handle each question with responses
	for _, q := range msg.Question {
		responses := handleQuestion(msg.Msg, q, this.Listener.Zone(), this.Services, this.Records)
		for _, response := range responses {
			// Ignore errors on send
			this.Listener.Send(response, msg.ifIndex)
		}
	}

	// Success
	return nil
}

func (this *Responder) Serve(ctx context.Context, r []gopi.ServiceRecord) error {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	// Set services which will be served
	if err := this.SetServices(r); err != nil {
		return err
	} else if s, _ := this.Services(); len(s) == 0 {
		return gopi.ErrBadParameter.WithPrefix("Serve")
	} else {
		this.Debug("Serve:", this.names)
	}

	// Wait for context to end
	<-ctx.Done()

	// TODO: Emit messages on listener with TTL=0

	return nil
}

func (this *Responder) NewServiceRecord(service string, name string, port uint16, txt []string, flags gopi.ServiceFlag) (gopi.ServiceRecord, error) {
	r := NewService(this.Listener.Zone())

	// Set service and name
	r.service = fqn(service) + r.zone
	r.name = fqn(Quote(name)) + r.service

	// Add host
	if host, err := os.Hostname(); err != nil {
		return nil, err
	} else if port == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("port")
	} else {
		host = fqn(host)
		if strings.HasSuffix(host, r.zone) == false {
			host = host + r.zone
		}
		r.host = append(r.host, target{host, port, 1})
	}

	// Addr
	if flags&gopi.SERVICE_FLAG_IP4 != 0 || flags == gopi.SERVICE_FLAG_NONE {
		r.a = this.Listener.AddrForIface(0, gopi.SERVICE_FLAG_IP4)
	}
	if flags&gopi.SERVICE_FLAG_IP6 != 0 || flags == gopi.SERVICE_FLAG_NONE {
		r.aaaa = this.Listener.AddrForIface(0, gopi.SERVICE_FLAG_IP6)
	}

	// Add txt
	r.txt = txt

	// Return success
	return r, nil
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func handleQuestion(msg *dns.Msg, question dns.Question, zone string, f1 FuncServices, f2 FuncRecordsForService) []*dns.Msg {
	// If forcing over unicast, ignore (RFC 6762, section 18.12)
	if question.Qclass&(1<<15) != 0 {
		return nil
	}

	// If in the wrong zone, then don't handle
	if strings.HasSuffix(question.Name, zone) == false {
		return nil
	}

	// Remove the zone from the question
	questionName := strings.TrimSuffix(question.Name, zone)
	switch {
	case questionName == fqn(queryServices):
		return handleEnum(msg, question, zone, f1)
	default:
		fmt.Println("TODO: Unhandled question:", questionName)
		return nil
	}
}

func handleEnum(req *dns.Msg, question dns.Question, zone string, fn FuncServices) []*dns.Msg {
	// Check incoming parameters
	if req == nil || fn == nil {
		return nil
	}
	// Handle PTR and ANY only
	if question.Qtype != dns.TypeANY && question.Qtype != dns.TypePTR {
		return nil
	}
	// Get services and the ttl
	services, ttl := fn()
	if len(services) == 0 {
		return nil
	}
	if ttl == 0 {
		ttl = queryDefaultTTL
	}

	// One message per service name
	msgs := make([]*dns.Msg, 0, len(services))
	for _, service := range services {
		rr := &dns.PTR{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    uint32(ttl),
			},
			Ptr: fqn(service) + zone,
		}
		msgs = append(msgs, prepareResponse(req, rr))
	}
	return msgs
}

func prepareResponse(req *dns.Msg, answers ...dns.RR) *dns.Msg {
	var queryId uint16
	if len(answers) == 0 {
		return nil
	}
	//if unicast { queryId = msg.Id }
	return &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id:            queryId,
			Response:      true,
			Opcode:        dns.OpcodeQuery,
			Authoritative: true,
		},
		Compress: true,
		Answer:   answers,
	}
}
