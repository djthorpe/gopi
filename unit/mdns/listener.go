/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	dns "github.com/miekg/dns"
	ipv4 "golang.org/x/net/ipv4"
	ipv6 "golang.org/x/net/ipv6"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Listener struct {
	domain   string
	ifaces   []net.Interface
	ip4      *ipv4.PacketConn
	ip6      *ipv6.PacketConn
	end      int32
	errors   chan<- error
	messages chan<- *dns.Msg

	sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MDNS_DEFAULT_DOMAIN = "local."
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	MDNS_ADDR_IPV4 = &net.UDPAddr{IP: net.ParseIP("224.0.0.251"), Port: 5353}
	MDNS_ADDR_IPV6 = &net.UDPAddr{IP: net.ParseIP("ff02::fb"), Port: 5353}
)

////////////////////////////////////////////////////////////////////////////////
// INIT / DESTROY

func (this *Listener) Init(config Discovery, errors chan<- error, messages chan<- *dns.Msg) error {
	if config.Domain == "" {
		config.Domain = MDNS_DEFAULT_DOMAIN
	}
	if ifaces, err := listMulticastInterfaces(config.Interface); err != nil {
		return err
	} else if errors == nil {
		return gopi.ErrBadParameter.WithPrefix("errors")
	} else if messages == nil {
		return gopi.ErrBadParameter.WithPrefix("messages")
	} else {
		this.domain = config.FQDomain()
		this.ifaces = ifaces
		this.end = 0
		this.errors = errors
		this.messages = messages
	}

	// Connect to interfaces
	if config.Flags&gopi.RPC_FLAG_INET_V4 != 0 {
		if ip4, err := joinUdp4Multicast(this.ifaces, MDNS_ADDR_IPV4); err != nil {
			return err
		} else {
			this.ip4 = ip4
		}
	}
	if config.Flags&gopi.RPC_FLAG_INET_V6 != 0 {
		if ip6, err := joinUdp6Multicast(this.ifaces, MDNS_ADDR_IPV6); err != nil {
			return err
		} else {
			this.ip6 = ip6
		}
	}

	// Start listening to connections
	go this.recv_loop4(this.ip4)
	go this.recv_loop6(this.ip6)

	// Success
	return nil
}

func (this *Listener) Destroy() error {
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
	this.errors = nil
	this.messages = nil

	// Return compound errors
	return errs.ErrorOrSelf()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this Listener) String() string {
	return "<Listener domain=" + strconv.Quote(this.domain) + ">"
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

// recv_loop is a long running routine to receive packets from an interface
func (this *Listener) recv_loop4(conn *ipv4.PacketConn) {
	// Sanity check
	if conn == nil {
		return
	}

	// Indicate end of loop
	this.Add(1)
	defer this.Done()

	// Perform loop
	buf := make([]byte, 65536)
	for atomic.LoadInt32(&this.end) == 0 {
		if n, cm, from, err := conn.ReadFrom(buf); err != nil {
			continue
		} else if cm == nil {
			continue
		} else if err := this.parse_packet(buf[:n], cm.IfIndex, from); err != nil {
			if this.errors != nil {
				this.errors <- err
			}
		}
	}
}

func (this *Listener) recv_loop6(conn *ipv6.PacketConn) {
	// Sanity check
	if conn == nil {
		return
	}

	// Indicate end of loop
	this.Add(1)
	defer this.Done()

	// Perform loop
	buf := make([]byte, 65536)
	for atomic.LoadInt32(&this.end) == 0 {
		if n, cm, from, err := conn.ReadFrom(buf); err != nil {
			continue
		} else if cm == nil {
			continue
		} else if err := this.parse_packet(buf[:n], cm.IfIndex, from); err != nil {
			if this.errors != nil {
				this.errors <- err
			}
		}
	}
}

// parse packets into service records
func (this *Listener) parse_packet(packet []byte, ifIndex int, from net.Addr) error {
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
	if this.messages != nil {
		this.messages <- &msg
	}
	return nil
}
