/*
Go Language Raspberry Pi Interface
(c) Copyright David Thorpe 2016
All Rights Reserved

For Licensing and Usage information, please see LICENSE.md
*/
package dispmanx

import (
	"errors"
)

var (
	ErrorDisplay = errors.New("Display Error")
	ErrorGetInfo = errors.New("GetInfo Error")
	ErrorUpdate = errors.New("Update Error")
	ErrorUpdateInProgress = errors.New("Update already in progress")
)
