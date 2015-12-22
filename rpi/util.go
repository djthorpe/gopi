/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

////////////////////////////////////////////////////////////////////////////////

type State struct {
}

////////////////////////////////////////////////////////////////////////////////

func New() *State {
	// create this object
	this := new(State)

	// initialize
	BCMHostInit()
	VCGenCmdInit()

	// Return this
	return this
}

func (this *State) Terminate() {
	VCGenCmdStop()
	BCMHostTerminate()
}
