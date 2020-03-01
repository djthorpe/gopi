/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	"context"
	"fmt"
	"strings"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	dns "github.com/miekg/dns"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Register struct {
	Listener ListenerIface
}

type register struct {
	listener ListenerIface
	stop     gopi.Channel
	records  map[string]gopi.RPCServiceRecord
	names    map[string]uint

	sync.Mutex
	sync.WaitGroup
	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Register) Name() string { return "gopi/mdns/register" }

func (config Register) New(log gopi.Logger) (gopi.Unit, error) {

	// Check parameters
	if config.Listener == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("Listener")
	}

	// Init
	this := new(register)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else {
		this.listener = config.Listener
		this.records = make(map[string]gopi.RPCServiceRecord)
		this.names = make(map[string]uint)
	}

	// Start background message processor
	if this.stop = this.listener.Subscribe(QUEUE_MESSAGES, this.EventHandler); this.stop == 0 {
		return nil, gopi.ErrInternalAppError
	}

	// Success
	return this, nil
}

func (this *register) Close() error {

	// Wait for EventHandler to end
	this.listener.Unsubscribe(this.stop)

	// Wait for register
	this.WaitGroup.Wait()

	// Release resources
	this.listener = nil
	this.records = nil
	this.names = nil

	// Return success
	return this.Unit.Close()
}

func (this *register) String() string {
	return fmt.Sprintf("<gopi.mDNS.Register %v services=%v>", this.listener, this.records)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.RPCServiceRegister

// Register service record, and de-register when deadline is exceeded
func (this *register) Register(ctx context.Context, record gopi.RPCServiceRecord) error {
	// Indicate Close() method should wait until all Register methods have ended
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	// Add record
	if record_, err := this.addRecord(record); err != nil {
		return err
	} else if err := this.sendRegister(record_); err != nil {
		return err
	} else {

		// Wait for stop or context end
		select {
		// TODO: Also react to close stop signal
		case <-ctx.Done():
			break
		}

		// Remove record
		if err := this.deleteRecord(record_); err != nil {
			return err
		} else if err := this.sendUnregister(record_); err != nil {
			return err
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// ADD AND REMOVE RECORDS

func (this *register) keyForRecord(record gopi.RPCServiceRecord) string {
	// Check parameters
	if record.Name == "" || record.Service == "" {
		return ""
	}
	if reService.MatchString(record.Service) == false {
		return ""
	}
	// Generate a unique key from the name and service
	return Quote(record.Name) + "." + record.Service + "." + this.listener.Zone()
}

func (this *register) addRecord(record gopi.RPCServiceRecord) (gopi.RPCServiceRecord, error) {
	this.Lock()
	defer this.Unlock()

	if key := this.keyForRecord(record); key == "" {
		return gopi.RPCServiceRecord{}, gopi.ErrBadParameter.WithPrefix("record")
	} else if _, exists := this.records[key]; exists {
		return gopi.RPCServiceRecord{}, gopi.ErrDuplicateItem.WithPrefix(key)
	} else {
		// Fully qualify hostname
		if strings.HasSuffix(record.Host, "."+this.listener.Zone()) == false {
			record.Host = record.Host + "." + this.listener.Zone()
		}
		// Set record
		this.records[key] = record
		if counter, exists := this.names[record.Service]; exists == false {
			this.names[record.Service] = 1
		} else {
			this.names[record.Service] = counter + 1
		}
		return record, nil
	}
}

func (this *register) deleteRecord(record gopi.RPCServiceRecord) error {
	this.Lock()
	defer this.Unlock()

	if key := this.keyForRecord(record); key == "" {
		return gopi.ErrBadParameter.WithPrefix("record")
	} else if _, exists := this.records[key]; exists == false {
		return gopi.ErrNotFound.WithPrefix(key)
	} else {
		delete(this.records, key)
		counter := this.names[record.Service] - 1
		if counter == 0 {
			delete(this.names, record.Service)
		} else {
			this.names[record.Service] = counter
		}
		return nil
	}
}

func (this *register) matchesServiceName(service string) bool {
	this.Lock()
	defer this.Unlock()

	_, exists := this.names[service]
	return exists
}

func (this *register) matchesInstanceName(key string) bool {
	this.Lock()
	defer this.Unlock()

	_, exists := this.records[key]
	return exists
}

func (this *register) matchesHostName(host string) bool {
	this.Lock()
	defer this.Unlock()

	// Fully qualify hostname
	if strings.HasSuffix(host, "."+this.listener.Zone()) == false {
		host = host + "." + this.listener.Zone()
	}
	// Scan records for hostname
	for _, record := range this.records {
		if record.Host == host {
			return true
		}
	}
	return false
}

func (this *register) recordsForServiceName(service string) []gopi.RPCServiceRecord {
	this.Lock()
	defer this.Unlock()

	records := make([]gopi.RPCServiceRecord, 0, len(this.records))
	for _, record := range this.records {
		if record.Service == service {
			records = append(records, record)
		}
	}
	return records
}

func (this *register) sendRegister(record gopi.RPCServiceRecord) error {
	if key := this.keyForRecord(record); key == "" {
		return gopi.ErrInternalAppError
	} else if msgs := this.handleEnum(dns.Question{
		Name:  DISCOVERY_SERVICE_QUERY + "." + this.listener.Zone(),
		Qtype: dns.TypeANY,
	}); msgs == nil {
		return gopi.ErrInternalAppError
	} else {
		msgs = append(msgs, this.handleRecord(dns.Question{
			Name:  record.Service + "." + this.listener.Zone(),
			Qtype: dns.TypeANY,
		}, record))
		for _, msg := range msgs {
			if err := this.listener.SendAll(msg); err != nil {
				return err
			}
		}
	}

	// Success
	return nil
}

func (this *register) sendUnregister(record gopi.RPCServiceRecord) error {
	if key := this.keyForRecord(record); key == "" {
		return gopi.ErrInternalAppError
	} else if msg := prepareResponse(&dns.PTR{
		Hdr: dns.RR_Header{
			Name:   record.Service + "." + this.listener.Zone(),
			Rrtype: dns.TypePTR,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		Ptr: key,
	}); msg == nil {
		return gopi.ErrInternalAppError
	} else if err := this.listener.SendAll(msg); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND PROCESSOR

func (this *register) EventHandler(value interface{}) {
	if msg, ok := value.(*dns.Msg); ok && msg == nil {
		return
	} else if len(msg.Question) == 0 {
		return
	} else if err := this.ProcessQuestion(msg); err != nil {
		this.Log.Error(err.(error))
	}
}

func (this *register) ProcessQuestion(msg *dns.Msg) error {
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
		if responses, err := this.handleQuestion(msg, q); err != nil {
			return err
		} else {
			for _, response := range responses {
				if err := this.listener.SendAll(response); err != nil {
					return err
				}
			}
		}
	}

	// Success
	return nil
}

func (this *register) handleQuestion(msg *dns.Msg, question dns.Question) ([]*dns.Msg, error) {
	// If forcing over unicast, ignore (RFC 6762, section 18.12)
	if question.Qclass&(1<<15) != 0 {
		return nil, nil
	}
	// If in the wrong zone, then don't handle
	if strings.HasSuffix(question.Name, "."+this.listener.Zone()) == false {
		return nil, nil
	}

	// Remove the zone from the question
	questionName := strings.TrimSuffix(question.Name, "."+this.listener.Zone())

	// Handle each question
	switch {
	case questionName == DISCOVERY_SERVICE_QUERY:
		return this.handleEnum(question), nil
	case this.matchesServiceName(questionName):
		return this.handleRecords(questionName, question), nil
	case this.matchesInstanceName(questionName):
		return this.handleInstance(questionName, question), nil
	case this.matchesHostName(questionName):
		if question.Qtype == dns.TypeA || question.Qtype == dns.TypeAAAA {
			return this.handleAddresses(question.Name, question), nil
		} else {
			fmt.Println("Unhandled DNS Qtype", question.Qtype)
		}
	}
	//return nil, fmt.Errorf("Unhandled question: %s", question.Name)
	return nil, nil
}

func (this *register) handleEnum(question dns.Question) []*dns.Msg {
	// Handle PTR and ANY
	if question.Qtype != dns.TypeANY && question.Qtype != dns.TypePTR {
		return nil
	}

	// One message per service name
	msgs := make([]*dns.Msg, 0, len(this.names))
	for k := range this.names {
		rr := &dns.PTR{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    uint32(MDNS_DEFAULT_TTL),
			},
			Ptr: k + "." + this.listener.Zone(),
		}
		msgs = append(msgs, prepareResponse(rr))
	}
	return msgs
}

func (this *register) handleRecord(question dns.Question, record gopi.RPCServiceRecord) *dns.Msg {
	key := this.keyForRecord(record)
	if key == "" {
		return nil
	}

	answers := []dns.RR{&dns.PTR{
		Hdr: dns.RR_Header{
			Name:   question.Name,
			Rrtype: dns.TypePTR,
			Class:  dns.ClassINET,
			Ttl:    uint32(MDNS_DEFAULT_TTL),
		},
		Ptr: key,
	}}

	// Append record answers
	for _, answer := range this.handleRecordAnswers(dns.Question{
		Name:  key,
		Qtype: dns.TypeANY,
	}, key, record) {
		if answer != nil {
			answers = append(answers, answer)
		}
	}

	return prepareResponse(answers...)
}

func (this *register) handleRecords(service string, question dns.Question) []*dns.Msg {
	// Handle PTR and ANY
	if question.Qtype != dns.TypeANY && question.Qtype != dns.TypePTR {
		return nil
	}
	// Get messages for each record
	records := this.recordsForServiceName(service)
	msgs := []*dns.Msg{}
	for _, record := range records {
		if msg := this.handleRecord(question, record); msg != nil {
			msgs = append(msgs, msg)
		}
	}
	return msgs
}

func (this *register) handleInstance(service string, question dns.Question) []*dns.Msg {
	if record, exists := this.records[service]; exists == false {
		return nil
	} else {
		answers := this.handleRecordAnswers(question, service, record)
		msg := prepareResponse(answers...)
		return []*dns.Msg{msg}
	}
}

func (this *register) handleAddresses(host string, question dns.Question) []*dns.Msg {
	msgs := make([]*dns.Msg, 0, len(this.records))
	for _, record := range this.records {
		if record.Host == host {
			if answers := this.handleRecordAnswers(question, host, record); len(answers) > 0 {
				msgs = append(msgs, prepareResponse(answers...))
			}
		}
	}
	return msgs
}

func (this *register) handleRecordAnswers(question dns.Question, instance string, record gopi.RPCServiceRecord) []dns.RR {
	switch question.Qtype {
	case dns.TypeANY:
		recs := this.handleRecordAnswers(dns.Question{
			Qtype: dns.TypeSRV,
			Name:  instance,
		}, instance, record)
		return append(recs, this.handleRecordAnswers(dns.Question{
			Qtype: dns.TypeTXT,
			Name:  instance,
		}, instance, record)...)
	case dns.TypeSRV:
		srv := &dns.SRV{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypeSRV,
				Class:  dns.ClassINET,
				Ttl:    uint32(MDNS_DEFAULT_TTL),
			},
			Priority: 10,
			Weight:   1,
			Port:     uint16(record.Port),
			Target:   record.Host,
		}
		// Add the A record
		recs := append([]dns.RR{srv}, this.handleRecordAnswers(dns.Question{
			Qtype: dns.TypeA,
			Name:  instance,
		}, instance, record)...)
		// Add the AAAA record
		return append(recs, this.handleRecordAnswers(dns.Question{
			Qtype: dns.TypeAAAA,
			Name:  instance,
		}, instance, record)...)
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
}

func prepareResponse(answers ...dns.RR) *dns.Msg {
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
