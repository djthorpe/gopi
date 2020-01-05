/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved

	Documentation https://gopi.mutablelogic.com/
	For Licensing and Usage information, please see LICENSE.md
*/

package event

import "github.com/djthorpe/gopi"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type nullEvent struct{}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	NullEvent = &nullEvent{}
)

////////////////////////////////////////////////////////////////////////////////
// EVENT INTERFACE IMPLEMENTATION

func (*nullEvent) Source() gopi.Driver {
	return nil
}

func (*nullEvent) Name() string {
	return "NullEvent"
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (*nullEvent) String() string {
	return "<gopi.NullEvent>{}"
}
