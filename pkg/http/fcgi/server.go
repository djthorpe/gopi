package fcgi

// This file implements FastCGI from the perspective of a child process.

import (
	"context"
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

	// Private variables to flag shutdown
	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc
}

func (s *Server) ListenAndServe() error {
	var wg sync.WaitGroup

	// Create listender
	if s.Path == "" {
		if l, err := net.FileListener(os.Stdin); err != nil {
			return err
		} else {
			s.listener = l
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

		if l, err := net.Listen("unix", s.Path); err != nil {
			return err
		} else {
			s.listener = l
		}
	}
	defer s.listener.Close()

	// Set default handler
	if s.Handler == nil {
		s.Handler = http.DefaultServeMux
	}

	// Set up semapore which when closed ends the loop
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// Continue accepting requests until shutdown
FOR_LOOP:
	for {
		select {
		case <-s.ctx.Done():
			break FOR_LOOP
		default:
			rw, err := s.listener.Accept()
			if err != nil {
				return err
			}
			c := newChild(rw, s.Handler)
			wg.Add(1)
			go func() {
				c.serve()
				wg.Done()
			}()
		}
	}

	// Wait until all connections served
	wg.Wait()

	// Return success
	return nil
}

func (s *Server) Close() error {
	if s.cancel != nil {
		s.listener.Close()
		s.cancel()
	}
	return nil
}
