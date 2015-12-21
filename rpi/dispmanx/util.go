/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package dispmanx

////////////////////////////////////////////////////////////////////////////////

type Display struct {
	device  uint32        // Device
	display DisplayHandle // Handle to the display
	update  UpdateHandle  // Update handle
	info    ModeInfo      // Information about the display
}

////////////////////////////////////////////////////////////////////////////////

func New(device uint32) (*Display, error) {
	// create display
	this := new(Display)
	this.device = device
	this.display = DisplayHandle(0)
	this.update = UpdateHandle(DISPMANX_NO_HANDLE)

	// Open display
	display, err := DisplayOpen(this.device)
	if err != nil {
		return nil, err
	} else {
		this.display = display
	}

	// Get information
	if err := DisplayGetInfo(this.display,&this.info); err != nil {
		DisplayClose(this.display)
		return nil,err
	}

	return this, nil
}

func (this *Display) Terminate() error {
	if this.display != DisplayHandle(DISPMANX_NO_HANDLE) {
		if err := DisplayClose(this.display); err != nil {
			return err
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (this *Display) StartUpdate(priority uint32) error {
	if this.update != UpdateHandle(DISPMANX_NO_HANDLE) {
		return ErrorUpdateInProgress
	}

	if this.update = UpdateStart(priority); this.update == UpdateHandle(DISPMANX_NO_HANDLE) {
		return ErrorUpdate
	}
	return nil
}

func (this *Display) EndUpdate() error {
	if this.update != UpdateHandle(DISPMANX_NO_HANDLE) {
		if UpdateSubmitSync(this.update) != DISPMANX_SUCCESS {
			return ErrorUpdate
		}
	}
	this.update = UpdateHandle(DISPMANX_NO_HANDLE)
	return nil
}

////////////////////////////////////////////////////////////////////////////////


