/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	"net"
	"regexp"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DISCOVERY_SERVICE_QUERY = "_services._dns-sd._udp"
	MDNS_DEFAULT_DOMAIN     = "local."
	MDNS_DEFAULT_TTL        = 120

	QUERY_REPEAT   = 2   // Number of times to repeat a message
	DELTA_QUERY_MS = 500 // Maximum time to wait between repeats

	// Pulisher queue numbers
	QUEUE_MESSAGES = 0
	QUEUE_ERRORS   = 1
	QUEUE_NAME     = 2
	QUEUE_RECORD   = 3
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	MULTICAST_ADDR_IPV4 = &net.UDPAddr{IP: net.ParseIP("224.0.0.251"), Port: 5353}
	MULTICAST_ADDR_IPV6 = &net.UDPAddr{IP: net.ParseIP("ff02::fb"), Port: 5353}
)

var (
	reService = regexp.MustCompile("^_[A-Za-z][A-Za-z0-9\\-]*\\._(tcp|udp)$")
)
