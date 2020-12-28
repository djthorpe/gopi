package http

import (
	"net"
	"os"
)

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func createListener(network, addr string) (net.Listener, error) {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		network = "unix"
		if _, err := os.Stat(addr); os.IsNotExist(err) == false {
			if err := os.Remove(addr); err != nil {
				return nil, err
			}
		}
	}
	if listener, err := net.Listen(network, addr); err != nil {
		return nil, err
	} else {
		return listener, nil
	}
}
