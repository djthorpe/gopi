// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package files

import (
	// Frameworks
	unix "golang.org/x/sys/unix"
)

type Pipe struct {
	fd [2]int // read[0] and write[1] file descriptors
}

func (this *Pipe) Init() error {
	if err := unix.Pipe2(this.fd[:], unix.O_NONBLOCK); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Pipe) Close() error {
	if this.fd[1] != -1 {
		if err := unix.Close(this.fd[1]); err != nil {
			return err
		}
	}
	if this.fd[0] != -1 {
		if err := unix.Close(this.fd[0]); err != nil {
			return err
		}
	}
	return nil
}

func (this *Pipe) ReadFd() uintptr {
	return uintptr(this.fd[0])
}


func (this *Pipe) WriteFd() uintptr {
	return uintptr(this.fd[1])
}


func (this *Pipe) Wake() error {
	buf := make([]byte, 1)
	if n, err := unix.Write(this.fd[1], buf); n == -1 {
		if err == unix.EAGAIN {
			return nil
		} else {
			return err
		}
	} else {
		return nil
	}
}


func (this *Pipe) Clear() error {
	buf := make([]byte, 100)
	FOR_LOOP:
	for {
		if n, err := unix.Read(this.fd[0], buf); n == -1 {
			if err == unix.EAGAIN {
				break FOR_LOOP
			} else {
				return err
			}
		} else if n == 0 {
			break FOR_LOOP
		}
	}
	return nil
}
