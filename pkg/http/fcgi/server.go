package fcgi

// This file implements FastCGI from the perspective of a child process.

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
)

// A Server defines parameters for running an FCGI server.
// The zero value for Server is a valid configuration, which responds
// to requests over stdin
type Server struct {
	// Path is the path to the socket, which is used to create the
	// listener. If path is empty, os.Stdin is used
	Path string

	// handler to invoke, http.DefaultServeMux if nil
	Handler http.Handler

	ch chan struct{}
}

func (s *Server) ListenAndServe() error {
	var l net.Listener
	var err error
	var wg sync.WaitGroup

	// Create listender
	if s.Path == "" {
		l, err = net.FileListener(os.Stdin)
		if err != nil {
			return err
		}
	} else {
		// Check for existing file and remove it. Cannot use a directory
		// as a socket
		if stat, err := os.Stat(s.Path); os.IsNotExist(err) {
			// File does not exist, so no nothing
		} else if err != nil {
			return err
		} else if stat.IsDir() {
			return fmt.Errorf("Cannot use an existing directory")
		} else if err := os.RemoveAll(s.Path); err != nil {
			return err
		}

		l, err = net.Listen("unix", s.Path)
		if err != nil {
			return err
		}
	}
	defer l.Close()

	// Set default handler
	if s.Handler == nil {
		s.Handler = http.DefaultServeMux
	}

	// Set up semapore which when closed ends the loop
	s.ch = make(chan struct{})

	// Continue accepting requests until shutdown
	for range s.ch {
		rw, err := l.Accept()
		if err != nil {
			return err
		}
		c := newChild(rw, s.Handler)
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.serve()
		}()
	}

	// Wait until all connections served
	wg.Wait()

	// Return success
	return nil
}

func (s *Server) Close() error {
	close(s.ch)
	return nil
}
