// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

type I2C struct {
	Bus uint
}

type i2c struct {
	log   gopi.Logger
	bus   uint
	slave uint8
	dev   *os.File
	funcs I2CFunction
	lock  sync.Mutex
}
