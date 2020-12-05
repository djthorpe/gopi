package mdns

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/djthorpe/gopi/v3"
	"github.com/miekg/dns"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Responder struct {
	gopi.Unit
	gopi.Publisher
	sync.WaitGroup
	gopi.Logger
	*Listener
}

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
			if msg, ok := evt.(*msgevent); ok {
				if err := this.ProcessQuestion(msg); err != nil {
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
		responses := handleQuestion(msg.Msg, q, this.Listener.Zone())
		for _, response := range responses {
			// Ignore errors on send
			this.Listener.Send(response, msg.ifIndex)
		}
	}

	// Success
	return nil
}

func (this *Responder) Serve(ctx context.Context, r []gopi.ServiceRecord) error {
	this.Print("TODO: Serve:", r)
	<-ctx.Done()
	return ctx.Err()
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func handleQuestion(msg *dns.Msg, question dns.Question, zone string) []*dns.Msg {
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
		return handleEnum(msg, question, zone)
	default:
		fmt.Println("TODO: Unhandled question:", questionName)
		return nil
	}
}

func handleEnum(req *dns.Msg, question dns.Question, zone string) []*dns.Msg {
	// Handle PTR and ANY
	if question.Qtype != dns.TypeANY && question.Qtype != dns.TypePTR {
		return nil
	}

	// TODO
	services := []string{
		"_gopi._tcp.",
	}

	// One message per service name
	msgs := make([]*dns.Msg, 0, len(services))
	for _, service := range services {
		rr := &dns.PTR{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    queryDefaultTTL,
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
