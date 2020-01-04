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
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	dns "github.com/miekg/dns"
	ipv4 "golang.org/x/net/ipv4"
	ipv6 "golang.org/x/net/ipv6"
)

////////////////////////////////////////////////////////////////////////////////
// Listener Interface

type ListenerIface interface {
	// Return properties
	Zone() string

	// Perform queries
	QueryAll(ctx context.Context, msg *dns.Msg, count uint) error

	// Send responses
	SendAll(msg *dns.Msg) error

	// Implements pub/sub interface
	gopi.Publisher
}

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Listener struct {
	Domain    string
	Interface net.Interface
	Flags     gopi.RPCFlag
}

type listener struct {
	domain string
	ifaces []net.Interface
	ip4    *ipv4.PacketConn
	ip6    *ipv6.PacketConn
	end    int32

	base.Unit
	base.Publisher
	sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Listener) Name() string { return "gopi.mDNS.Listener" }

func (config Listener) FQDomain() string {
	if config.Domain == "" {
		return ""
	} else {
		return strings.Trim(config.Domain, ".") + "."
	}
}

func (config Listener) New(log gopi.Logger) (gopi.Unit, error) {
	// Check parameters
	if config.Domain == "" {
		config.Domain = MDNS_DEFAULT_DOMAIN
	}
	if config.Flags&gopi.RPC_FLAG_INET_V4 == 0 && config.Flags&gopi.RPC_FLAG_INET_V6 == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("Flags")
	}

	// Create listener
	this := new(listener)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else if ifaces, err := listMulticastInterfaces(config.Interface); err != nil {
		return nil, err
	} else {
		this.domain = config.FQDomain()
		this.ifaces = ifaces
		this.end = 0
	}

	// Connect to interfaces
	if config.Flags&gopi.RPC_FLAG_INET_V4 != 0 {
		if ip4, err := joinUdp4Multicast(this.ifaces, MULTICAST_ADDR_IPV4); err != nil {
			return nil, err
		} else {
			this.ip4 = ip4
		}
	}
	if config.Flags&gopi.RPC_FLAG_INET_V6 != 0 {
		if ip6, err := joinUdp6Multicast(this.ifaces, MULTICAST_ADDR_IPV6); err != nil {
			return nil, err
		} else {
			this.ip6 = ip6
		}
	}

	// Start listening to connections
	go this.recv_loop4(this.ip4)
	go this.recv_loop6(this.ip6)

	// Return success
	return this, nil
}

func (this *listener) Close() error {
	errs := gopi.NewCompoundError()

	// Indicate shutdown
	if !atomic.CompareAndSwapInt32(&this.end, 0, 1) {
		return nil
	}

	// Close connections
	if this.ip4 != nil {
		if err := this.ip4.Close(); err != nil {
			errs.Add(err)
		}
	}

	if this.ip6 != nil {
		if err := this.ip6.Close(); err != nil {
			errs.Add(err)
		}
	}

	// Wait for recv_loop go routines to end
	this.Wait()

	// Release resources
	this.ip4 = nil
	this.ip6 = nil
	this.ifaces = nil
	this.end = 0

	// Close units/publisher
	errs.Add(this.Publisher.Close())
	errs.Add(this.Unit.Close())

	// Return compound errors
	return errs.ErrorOrSelf()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *listener) String() string {
	if this.Closed {
		return "<gopi.Listener closed=true>"
	} else {
		return "<gopi.Listener zone=" + strconv.Quote(this.domain) + ">"
	}
}

////////////////////////////////////////////////////////////////////////////////
// RETURN PROPERTIES

func (this *listener) Zone() string {
	return this.domain
}

////////////////////////////////////////////////////////////////////////////////
// QUERY AND SEND

// QueryAll sends a message to all multicast addresses
func (this *listener) QueryAll(ctx context.Context, msg *dns.Msg, count uint) error {
	// Send out message a certain number of times
	ticker := time.NewTimer(1 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			if err := this.multicast_send(msg, 0); err != nil {
				return err
			}
			if count > 0 {
				// Restart timer to send query again
				r := time.Duration(rand.Intn(DELTA_QUERY_MS))
				ticker.Reset(time.Millisecond * r)
				count--
			}
		case <-ctx.Done():
			ticker.Stop()
			return ctx.Err()
		}
	}
}

// SendAll sends message to all multicast addresses and returns
func (this *listener) SendAll(msg *dns.Msg) error {
	return this.multicast_send(msg, 0)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func joinUdp6Multicast(ifaces []net.Interface, addr *net.UDPAddr) (*ipv6.PacketConn, error) {
	if len(ifaces) == 0 {
		return nil, gopi.ErrBadParameter
	} else if conn, err := net.ListenUDP("udp6", addr); err != nil {
		return nil, err
	} else if packet_conn := ipv6.NewPacketConn(conn); packet_conn == nil {
		return nil, conn.Close()
	} else {
		packet_conn.SetControlMessage(ipv6.FlagInterface, true)
		errs := gopi.NewCompoundError()
		for _, iface := range ifaces {
			if err := packet_conn.JoinGroup(&iface, &net.UDPAddr{IP: addr.IP}); err != nil {
				errs.Add(fmt.Errorf("JoinGroup6: %v: %v", iface.Name, err))
			}
		}
		if errs.ErrorOrSelf() == nil {
			return packet_conn, nil
		}
		errs.Add(conn.Close())
		return nil, errs.ErrorOrSelf()
	}
}

func joinUdp4Multicast(ifaces []net.Interface, addr *net.UDPAddr) (*ipv4.PacketConn, error) {
	if len(ifaces) == 0 {
		return nil, gopi.ErrBadParameter
	} else if conn, err := net.ListenUDP("udp4", addr); err != nil {
		return nil, err
	} else if packet_conn := ipv4.NewPacketConn(conn); packet_conn == nil {
		return nil, conn.Close()
	} else {
		packet_conn.SetControlMessage(ipv4.FlagInterface, true)
		errs := gopi.NewCompoundError()
		for _, iface := range ifaces {
			if err := packet_conn.JoinGroup(&iface, &net.UDPAddr{IP: addr.IP}); err != nil {
				if err_, ok := err.(*os.SyscallError); ok && err_.Err == syscall.EAFNOSUPPORT {
					continue
				} else {
					errs.Add(fmt.Errorf("JoinGroup4: %v: %v", iface.Name, err))
				}
			}
		}
		if errs.ErrorOrSelf() == nil {
			return packet_conn, nil
		}
		errs.Add(conn.Close())
		return nil, errs.ErrorOrSelf()
	}
}

func listMulticastInterfaces(iface net.Interface) ([]net.Interface, error) {
	if iface.Name != "" {
		if (iface.Flags&net.FlagUp) > 0 && (iface.Flags&net.FlagMulticast) > 0 {
			return []net.Interface{iface}, nil
		} else {
			return nil, fmt.Errorf("Interface %v is not up and/or multicast-enabled", iface.Name)
		}
	}
	if ifaces, err := net.Interfaces(); err != nil {
		return nil, err
	} else {
		interfaces := make([]net.Interface, 0, len(ifaces))
		for _, ifi := range ifaces {
			if (ifi.Flags & net.FlagUp) == 0 {
				continue
			}
			if (ifi.Flags & net.FlagMulticast) == 0 {
				continue
			}
			if addrs, err := ifi.MulticastAddrs(); err != nil || len(addrs) == 0 {
				continue
			}
			interfaces = append(interfaces, ifi)
		}
		if len(interfaces) > 0 {
			return interfaces, nil
		} else {
			return nil, fmt.Errorf("No multicast-enabled interface found")
		}
	}
}

// recv_loop4 is a long running routine to receive packets from an interface
func (this *listener) recv_loop4(conn *ipv4.PacketConn) {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	// Sanity check
	if conn == nil {
		return
	}

	// Perform loop
	buf := make([]byte, 65536)
	for atomic.LoadInt32(&this.end) == 0 {
		if n, cm, from, err := conn.ReadFrom(buf); err != nil {
			continue
		} else if cm == nil {
			continue
		} else if err := this.parse_packet(buf[:n], cm.IfIndex, from); err != nil {
			this.Publisher.Emit(QUEUE_ERRORS, err)
		}
	}
}

// recv_loop6 is a long running routine to receive packets from an interface
func (this *listener) recv_loop6(conn *ipv6.PacketConn) {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	// Sanity check
	if conn == nil {
		return
	}

	// Perform loop
	buf := make([]byte, 65536)
	for atomic.LoadInt32(&this.end) == 0 {
		if n, cm, from, err := conn.ReadFrom(buf); err != nil {
			continue
		} else if cm == nil {
			continue
		} else if err := this.parse_packet(buf[:n], cm.IfIndex, from); err != nil {
			this.Publisher.Emit(QUEUE_ERRORS, err)
		}
	}
}

// parse packets into service records
func (this *listener) parse_packet(packet []byte, ifIndex int, from net.Addr) error {
	var msg dns.Msg
	if err := msg.Unpack(packet); err != nil {
		return err
	}
	if msg.Opcode != dns.OpcodeQuery {
		return fmt.Errorf("Query with invalid Opcode %v (expected %v)", msg.Opcode, dns.OpcodeQuery)
	}
	if msg.Rcode != 0 {
		return fmt.Errorf("Query with non-zero Rcode %v", msg.Rcode)
	}
	if msg.Truncated {
		return fmt.Errorf("Support for DNS requests with high truncated bit not implemented")
	}
	this.Publisher.Emit(QUEUE_MESSAGES, &msg)
	return nil
}

// multicastSend sends a multicast response packet to a particular interface
// or all interfaces if 0
func (this *listener) multicast_send(msg *dns.Msg, ifIndex int) error {
	var buf []byte
	if msg == nil {
		return gopi.ErrBadParameter.WithPrefix("msg")
	} else if buf_, err := msg.Pack(); err != nil {
		return err
	} else {
		buf = buf_
	}
	if this.ip4 != nil {
		var cm ipv4.ControlMessage
		if ifIndex != 0 {
			cm.IfIndex = ifIndex
			this.ip4.WriteTo(buf, &cm, MULTICAST_ADDR_IPV4)
		} else {
			for _, intf := range this.ifaces {
				cm.IfIndex = intf.Index
				this.ip4.WriteTo(buf, &cm, MULTICAST_ADDR_IPV4)
			}
		}
	}
	if this.ip6 != nil {
		var cm ipv6.ControlMessage
		if ifIndex != 0 {
			cm.IfIndex = ifIndex
			this.ip6.WriteTo(buf, &cm, MULTICAST_ADDR_IPV6)
		} else {
			for _, intf := range this.ifaces {
				cm.IfIndex = intf.Index
				this.ip6.WriteTo(buf, &cm, MULTICAST_ADDR_IPV6)
			}
		}
	}
	// Success
	return nil
}
