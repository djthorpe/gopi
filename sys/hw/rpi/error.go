// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"errors"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	ErrGeneralCommand     = errors.New("General Command Error")
	ErrUnexpectedResponse = errors.New("Unexpected response")
)
