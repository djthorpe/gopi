package http

import (
	"net"
)

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func getFreePort() (int, error) {
	if addr, err := net.ResolveTCPAddr("tcp", ":0"); err != nil {
		return 0, err
	} else if listener, err := net.ListenTCP("tcp", addr); err != nil {
		return 0, err
	} else {
		defer listener.Close()
		return listener.Addr().(*net.TCPAddr).Port, nil
	}
}
