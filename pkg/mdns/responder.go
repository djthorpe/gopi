package mdns

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
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
		serviceKey := fqn(record.Service())
		instanceKey := fqn(record.Instance())
		if _, exists := this.records[serviceKey]; exists == false {
			this.names = append(this.names, serviceKey)
		}
		this.records[serviceKey] = append(this.records[serviceKey], record)
		this.records[instanceKey] = append(this.records[instanceKey], record)
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
		return gopi.ErrBadParameter.WithPrefix("Query with non-zero opcode")
	}
	if msg.Rcode != 0 {
		return gopi.ErrBadParameter.WithPrefix("Query with non-zero Rcode")
	}
	if msg.Truncated {
		return gopi.ErrBadParameter.WithPrefix("Query with high truncated bit")
	}

	// Handle each question with responses
	zone := this.Zone()
	for _, q := range msg.Question {
		// Only answer questions for this zone
		if strings.HasSuffix(q.Name, zone) == false {
			continue
		}

		// Process responses to question
		answers := handleQuestion(msg.Msg, q, zone, this.Services, this.Records)

		// Send answers, ignoring any errors
		this.SendAnswers(msg.ifIndex, answers)

		// Report on any unhandled questions through debugging
		if len(answers) == 0 && this.isRelevantQuestion(q, zone) {
			this.Debug("Unhandled question: ", strconv.Quote(q.Name), " of type: ", qTypeString(q.Qtype))
		}
	}

	// Success
	return nil
}

func (this *Responder) SendAnswers(ifIndex int, msgs []*dns.Msg) error {
	var result error
	for _, msg := range msgs {
		if err := this.Listener.Send(msg, ifIndex); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
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
		this.Debug("Serve: ", this.names)
	}

	// Register services and service records
	msgs := []*dns.Msg{}
	zone := this.Listener.Zone()
	msgs = append(msgs, answerEnum(dns.Question{
		Name:  fqn(queryServices) + zone,
		Qtype: dns.TypeANY,
	}, this.names, zone)...)
	for _, record := range r {
		msgs = append(msgs, answerServiceRecords(dns.Question{
			Name:  record.Service() + record.Zone(),
			Qtype: dns.TypeANY,
		}, []gopi.ServiceRecord{record}, queryDefaultTTL)...)
	}
	this.SendAnswers(0, msgs)

	// Wait for context to end
	<-ctx.Done()

	// De-register service records
	for _, record := range r {
		msg := prepareResponse(&dns.PTR{
			Hdr: dns.RR_Header{
				Name:   record.Service() + record.Zone(),
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			Ptr: record.Instance() + record.Zone(),
		})
		this.SendAnswers(0, []*dns.Msg{msg})
	}

	// Set empty services
	if err := this.SetServices(nil); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Responder) NewServiceRecord(service string, name string, port uint16, txt []string, flags gopi.ServiceFlag) (gopi.ServiceRecord, error) {
	// Create service record
	r := NewService(this.Listener.Zone())

	// Set service and name
	r.service = fqn(service) + r.zone
	r.name = fqn(Quote(name)) + r.service

	// Add host
	if host, err := os.Hostname(); err != nil {
		return nil, err
	} else {
		host = fqn(host)
		if strings.HasSuffix(host, r.zone) == false {
			host = host + r.zone
		}
		r.host = host
		r.port = port
	}

	// Add addresses
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

// isRelevantQuestion returns true if a question has a suffix of a recorded
// service, ie, it's relevant to be answered
func (this *Responder) isRelevantQuestion(q dns.Question, zone string) bool {
	name := strings.TrimSuffix(q.Name, zone)
	for k := range this.records {
		if strings.HasSuffix(name, k) {
			return true
		}
	}
	return false
}

// handleQuestion prepares a response to be sent
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

	// Check for service, record or instance
	if questionName == fqn(queryServices) {
		return answerEnum(question, f1(), zone)
	} else if r := f2(questionName); len(r) > 0 {
		return answerServiceRecords(question, r, queryDefaultTTL)
	} else {
		return nil
	}
}

func answerEnum(question dns.Question, services []string, zone string) []*dns.Msg {
	// Check incoming parameters
	if len(services) == 0 {
		return nil
	}

	// Handle PTR and ANY only
	if question.Qtype != dns.TypeANY && question.Qtype != dns.TypePTR {
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
		msgs = append(msgs, prepareResponse(rr))
	}
	return msgs
}

func answerServiceRecords(question dns.Question, recs []gopi.ServiceRecord, ttl uint32) []*dns.Msg {
	// Check incoming parameters
	if len(recs) == 0 {
		return nil
	}

	// Get messages for each record
	msgs := []*dns.Msg{}
	for _, rec := range recs {
		if msg := answerRecord(question, rec, ttl); msg != nil {
			msgs = append(msgs, msg)
		}
	}
	return msgs
}

func answerRecord(question dns.Question, record gopi.ServiceRecord, ttl uint32) *dns.Msg {
	// Header
	answers := []dns.RR{&dns.PTR{
		Hdr: dns.RR_Header{
			Name:   question.Name,
			Rrtype: dns.TypePTR,
			Class:  dns.ClassINET,
			Ttl:    ttl,
		},
		Ptr: record.Instance() + record.Zone(),
	}}

	if question.Qtype == dns.TypePTR || question.Qtype == dns.TypeANY {
		answers = append(answers, answerSRV(question, record, ttl))
		answers = append(answers, answerA(question, record, ttl)...)
		answers = append(answers, answerAAAA(question, record, ttl)...)
		answers = append(answers, answerTxt(question, record, ttl))
	} else if question.Qtype == dns.TypeTXT {
		answers = append(answers, answerTxt(question, record, ttl))
	}

	/*
		fmt.Println("Question", question.Name, qTypeString(question.Qtype))
		for i, answer := range answers {
			fmt.Println("  ", i, answer)
		}
	*/

	return prepareResponse(answers...)
}

func answerSRV(question dns.Question, record gopi.ServiceRecord, ttl uint32) dns.RR {
	return &dns.SRV{
		Hdr: dns.RR_Header{
			Name:   record.Instance() + record.Zone(),
			Rrtype: dns.TypeSRV,
			Class:  dns.ClassINET,
			Ttl:    ttl,
		},
		Priority: 10,
		Weight:   1,
		Port:     record.Port(),
		Target:   record.Host(),
	}
}

func answerA(question dns.Question, record gopi.ServiceRecord, ttl uint32) []dns.RR {
	answers := []dns.RR{}
	for _, ip := range record.Addrs() {
		ip4 := ip.To4()
		if ip4 == nil {
			continue
		}
		answers = append(answers, &dns.A{
			Hdr: dns.RR_Header{
				Name:   record.Host(),
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    ttl,
			},
			A: ip4,
		})
	}
	return answers
}

func answerAAAA(question dns.Question, record gopi.ServiceRecord, ttl uint32) []dns.RR {
	answers := []dns.RR{}
	for _, ip := range record.Addrs() {
		if ip.To4() != nil {
			continue
		}
		ip6 := ip.To16()
		if ip6 == nil {
			continue
		}
		answers = append(answers, &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   record.Host(),
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    ttl,
			},
			AAAA: ip6,
		})
	}
	return answers
}

func answerTxt(question dns.Question, record gopi.ServiceRecord, ttl uint32) dns.RR {
	return &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   question.Name,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    ttl,
		},
		Txt: record.Txt(),
	}
}

func prepareResponse(answers ...dns.RR) *dns.Msg {
	if len(answers) == 0 {
		return nil
	}

	var queryId uint16
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
