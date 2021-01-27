package ts

import (
	"fmt"
	"io"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type NITSection struct {
	*Header
	*Table
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewNITSection(r io.Reader) (*NITSection, error) {
	this := new(NITSection)

	if header, err := NewHeader(r); err != nil {
		return nil, err
	} else if rows, err := NewTable(r); err != nil {
		return nil, err
	} else {
		this.Header = header
		this.Table = rows
	}

	// Return success
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *NITSection) String() string {
	str := "<nit"
	str += fmt.Sprint(" ", this.Header)
	str += fmt.Sprint(" ", this.Table)
	return str + ">"
}
