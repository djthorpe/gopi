/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

func init() {
	// Set gopi.NullEvent
	NullEvent = &nullevent{}
}

////////////////////////////////////////////////////////////////////////////////
// NULL EVENT IMPLEMENTATION

type nullevent struct{}

func (nullevent) Name() string       { return "gopi.NullEvent" }
func (nullevent) Source() Unit       { return nil }
func (nullevent) NS() EventNS        { return EVENT_NS_DEFAULT }
func (nullevent) Value() interface{} { return nil }
func (nullevent) String() string     { return "<gopi.NullEvent>" }
