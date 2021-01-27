package ts

import (
	"fmt"
	"io"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type NITSection struct {
	Header
	DTable
	STable
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (n *NITSection) Read(r io.Reader, length int) error {
	if err := n.Header.Read(r); err != nil {
		return err
	} else if err := n.DTable.Read(r); err != nil {
		return err
	} else if err := n.STable.Read(r); err != nil {
		return err
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (n NITSection) String() string {
	str := "<nit"
	str += fmt.Sprint(" ", n.Header)
	str += fmt.Sprint(" ", n.DTable)
	str += fmt.Sprint(" ", n.STable)
	return str + ">"
}
