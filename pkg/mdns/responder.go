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

// FuncServices returns fully-qualified service names
type FuncServices func() []string

// FuncRecordsForService returns service records for named service
type FuncRecordsForService func(string) []gopi.ServiceRecord

///////////////////////////////////////////////////////////////////////////////
// RUN

func (this *Responder) Run(ctx context.Context) error {
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
			if s := this.Services(); len(s) == 0 {
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

func (this *Responder) Services() []string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Return the current service names with unset TTL
	return this.names
}

func (this *Responder) Records(name string) []gopi.ServiceRecord {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	key := fqn(name)
	if r, exists := this.records[key]; exists {
		return r
	} else {
		return nil
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
	} else if s := this.Services(); len(s) == 0 {
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
	case len(f2(questionName)) > 0:
		return handleServiceRecords(msg, question, f2(questionName))
		return nil
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
	services := fn()
	if len(services) == 0 {
		return nil
	}

	// One message per service name
	msgs := make([]*dns.Msg, 0, len(services))
	for _, service := range services {
		rr := &dns.PTR{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    uint32(queryDefaultTTL),
			},
			Ptr: fqn(service) + zone,
		}
		msgs = append(msgs, prepareResponse(req, rr))
	}
	return msgs
}

func handleServiceRecords(req *dns.Msg, question dns.Question, recs []gopi.ServiceRecord) []*dns.Msg {
	// Check incoming parameters
	if len(recs) == 0 {
		return nil
	}
	// Handle PTR and ANY
	if question.Qtype != dns.TypeANY && question.Qtype != dns.TypePTR {
		return nil
	}

	// Get messages for each record
	msgs := []*dns.Msg{}
	for _, rec := range recs {
		if msg := handleRecord(req, question, rec); msg != nil {
			msgs = append(msgs, msg)
		}
	}
	return msgs
}

func handleRecord(req *dns.Msg, question dns.Question, record gopi.ServiceRecord) *dns.Msg {
	fmt.Println("ptr=", record.Ptr())
	// Header
	answers := []dns.RR{&dns.PTR{
		Hdr: dns.RR_Header{
			Name:   question.Name,
			Rrtype: dns.TypePTR,
			Class:  dns.ClassINET,
			Ttl:    queryDefaultTTL,
		},
		Ptr: record.Ptr(),
	}}

	// TODO Append record answers

	return prepareResponse(req, answers...)
}

func handleRecordAnswers(req *dns.Msg, question dns.Question, record gopi.ServiceRecord) []dns.RR {
	return nil
	/*
		switch question.Qtype {
		case dns.TypeANY:
			recs := handleRecordAnswers(req, dns.Question{
				Qtype: dns.TypeSRV,
				Name:  "i",
			}, record)
			return append(recs, handleRecordAnswers(req, dns.Question{
				Qtype: dns.TypeTXT,
				Name:  "i",
			}, record)...)
		case dns.TypeSRV:
			srv := &dns.SRV{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeSRV,
					Class:  dns.ClassINET,
					Ttl:    queryDefaultTTL,
				},
				Priority: 10,
				Weight:   1,
				Port:     uint16(record.Port),
				Target:   record.Host,
			}
			// Add the A record
			recs := append([]dns.RR{srv}, handleRecordAnswers(req, dns.Question{
				Qtype: dns.TypeA,
				Name:  "i",
			}, record)...)
			// Add the AAAA record
			return append(recs, handleRecordAnswers(req, dns.Question{
				Qtype: dns.TypeAAAA,
				Name:  "i",
			}, record)...)
		case dns.TypeA:
			var rr []dns.RR
			for _, ip := range record.Addrs {
				if ip4 := ip.To4(); ip4 != nil {
					rr = append(rr, &dns.A{
						Hdr: dns.RR_Header{
							Name:   record.Host,
							Rrtype: dns.TypeA,
							Class:  dns.ClassINET,
							Ttl:    uint32(MDNS_DEFAULT_TTL),
						},
						A: ip4,
					})
				}
			}
			return rr
		case dns.TypeAAAA:
			var rr []dns.RR
			for _, ip := range record.Addrs {
				if ip6 := ip.To16(); ip6 != nil {
					rr = append(rr, &dns.AAAA{
						Hdr: dns.RR_Header{
							Name:   record.Host,
							Rrtype: dns.TypeAAAA,
							Class:  dns.ClassINET,
							Ttl:    uint32(MDNS_DEFAULT_TTL),
						},
						AAAA: ip6,
					})
				}
			}
			return rr
		case dns.TypeTXT:
			txt := &dns.TXT{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    uint32(MDNS_DEFAULT_TTL),
				},
				Txt: record.Txt,
			}
			return []dns.RR{txt}
		}

		// Return nil
		return nil
	*/
}

func prepareResponse(req *dns.Msg, answers ...dns.RR) *dns.Msg {
	var queryId uint16
	if len(answers) == 0 {
		return nil
	}
	// if unicast { queryId = msg.Id }
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
