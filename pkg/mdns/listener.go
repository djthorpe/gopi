package mdns

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
	"github.com/miekg/dns"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Listener struct {
	sync.RWMutex
	sync.WaitGroup
	gopi.Unit
	gopi.Logger
	gopi.Publisher

	// Arguments
	domain, iface *string

	// Interfaces for listener
	ifaces []net.Interface

	// Bound listeners
	ip4 *ipv4.PacketConn
	ip6 *ipv6.PacketConn
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

const (
	emitRetryCount    = 3
	emitRetryDuration = 100 * time.Millisecond
)

var (
	MULTICAST_ADDR_IPV4 = &net.UDPAddr{IP: net.ParseIP("224.0.0.251"), Port: 5353}
	MULTICAST_ADDR_IPV6 = &net.UDPAddr{IP: net.ParseIP("ff02::fb"), Port: 5353}
)

///////////////////////////////////////////////////////////////////////////////
// INIT

func (this *Listener) Define(cfg gopi.Config) error {
	this.domain = cfg.FlagString("mdns.domain", "local.", "mDNS domain")
	this.iface = cfg.FlagString("mdns.iface", "", "mDNS listening interface")
	return nil
}

func (this *Listener) New(cfg gopi.Config) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Fully qualify domain (remove dots and add one to end)
	if *this.domain = strings.Trim(*this.domain, ".") + "."; *this.domain == "." {
		return gopi.ErrBadParameter.WithPrefix("-mdns.domain")
	}

	// Obtain the interfaces for listening
	if iface, err := interfaceForName(*this.iface); err != nil {
		return err
	} else if ifaces, err := multicastInterfaces(iface); err != nil {
		return err
	} else if len(ifaces) == 0 {
		return fmt.Errorf("No interfaces defined for listening")
	} else {
		this.ifaces = ifaces
	}

	// Join IP4
	if ip4, err := bindUdp4(this.ifaces, MULTICAST_ADDR_IPV4); err != nil {
		return err
	} else {
		this.ip4 = ip4
	}

	// Join IP6
	if ip6, err := bindUdp6(this.ifaces, MULTICAST_ADDR_IPV6); err != nil {
		return err
	} else {
		this.ip6 = ip6
	}

	// Return success
	return nil
}

func (this *Listener) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Close connections
	if this.ip4 != nil {
		if err := this.ip4.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.ip6 != nil {
		if err := this.ip6.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Wait until receive loops have completed
	this.WaitGroup.Wait()

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// RUN

func (this *Listener) Run(ctx context.Context) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Check to make sure there is  punlisher for emitting messages
	if this.Publisher == nil {
		return gopi.ErrInternalAppError
	}

	// Run4
	if this.ip4 != nil {
		this.WaitGroup.Add(1)
		go this.run4(ctx, this.ip4)
	}

	// Run6
	if this.ip6 != nil {
		this.WaitGroup.Add(1)
		go this.run6(ctx, this.ip6)
	}

	// Wait for cancels
	<-ctx.Done()

	// Return success
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Send a DNS message to a particular interface or all interfaces if 0
func (this *Listener) Send(msg *dns.Msg, ifIndex int) error {
	var buf []byte
	var result error

	if msg == nil {
		return gopi.ErrBadParameter.WithPrefix("msg")
	} else if buf_, err := msg.Pack(); err != nil {
		return err
	} else {
		buf = buf_
	}

	for i, q := range msg.Question {
		this.Debug("  ", i, " Send: ", q.Name, " type=", qTypeString(q.Qtype), " ifIndex=", ifIndex)
	}

	if this.ip4 != nil {
		var cm ipv4.ControlMessage
		if ifIndex != 0 {
			cm.IfIndex = ifIndex
			if _, err := this.ip4.WriteTo(buf, &cm, MULTICAST_ADDR_IPV4); err != nil {
				result = multierror.Append(result, err)
			}
		} else {
			for _, intf := range this.ifaces {
				cm.IfIndex = intf.Index
				if intf.Flags&net.FlagUp != 0 {
					this.ip4.WriteTo(buf, &cm, MULTICAST_ADDR_IPV4)
				}
			}
		}
	}

	if this.ip6 != nil {
		var cm ipv6.ControlMessage
		if ifIndex != 0 {
			cm.IfIndex = ifIndex
			if _, err := this.ip6.WriteTo(buf, &cm, MULTICAST_ADDR_IPV6); err != nil {
				result = multierror.Append(result, err)
			}
		} else {
			for _, intf := range this.ifaces {
				cm.IfIndex = intf.Index
				if intf.Flags&net.FlagUp != 0 {
					this.ip6.WriteTo(buf, &cm, MULTICAST_ADDR_IPV6)
				}
			}
		}
	}

	// Success
	return result
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Listener) run4(ctx context.Context, conn *ipv4.PacketConn) {
	defer this.WaitGroup.Done()

	buf := make([]byte, 65536)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if n, cm, from, err := conn.ReadFrom(buf); err != nil {
				continue
			} else if cm == nil {
				continue
			} else if msg, err := parseDnsPacket(buf[:n]); err != nil {
				this.Print("DNS Error:", err)
			} else if err := this.Publisher.Emit(NewMsgEvent(msg, from, cm.IfIndex), true); err != nil {
				this.Print("Emit Error:", err)
			}
		}
	}
}

func (this *Listener) run6(ctx context.Context, conn *ipv6.PacketConn) {
	defer this.WaitGroup.Done()

	buf := make([]byte, 65536)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if n, cm, from, err := conn.ReadFrom(buf); err != nil {
				continue
			} else if cm == nil {
				continue
			} else if msg, err := parseDnsPacket(buf[:n]); err != nil {
				this.Print("DNS Error:", err)
			} else if err := this.Publisher.Emit(NewMsgEvent(msg, from, cm.IfIndex), true); err != nil {
				this.Print("Emit Error:", err)
			}
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Listener) Zone() string {
	return *this.domain
}

func (this *Listener) AddrForIface(ifIndex int, flags gopi.ServiceFlag) []net.IP {
	ips := []net.IP{}
	for _, iface := range this.ifaces {
		if ifIndex != 0 && ifIndex != iface.Index {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ip, _, err := net.ParseCIDR(addr.String()); err != nil {
				this.Debug("AddrForIface: Error: ", addr.String())
				continue
			} else if ip.Equal(net.IPv6loopback) {
				continue
			} else if ip.String() == "127.0.0.1" {
				continue
			} else if flags == gopi.SERVICE_FLAG_IP6 && ip.To16() != nil {
				ips = append(ips, ip)
			} else if flags == gopi.SERVICE_FLAG_IP4 && ip.To4() != nil {
				ips = append(ips, ip)
			}
		}
	}
	return ips
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Listener) String() string {
	str := "<listener"
	str += fmt.Sprintf(" domain=%q", *this.domain)
	str += fmt.Sprintf(" ifaces=")
	for i, iface := range this.ifaces {
		if i > 0 {
			str += ","
		}
		str += iface.Name
	}

	return str + ">"
}
