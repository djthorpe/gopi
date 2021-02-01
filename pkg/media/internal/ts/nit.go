package ts

import (
	"fmt"
	"io"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type NITSection struct {
	Header
	DTable
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (n *NITSection) Read(r io.Reader, length int) error {
	return gopi.ErrNotImplemented
	/*
		if err := n.Header.Read(r); err != nil {
			return err
		} else if err := n.DTable.Read(r); err != nil {
			return err
		} else if err := n.STable.Read(r); err != nil {
			return err
		}

		// Return success
		return nil*/
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (n NITSection) String() string {
	str := "<nit"
	str += fmt.Sprint(" ", n.Header)
	str += fmt.Sprint(" ", n.DTable)
	return str + ">"
}
