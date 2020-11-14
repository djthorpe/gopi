package ping

import (
	gopi "github.com/djthorpe/gopi/v3"
)

type stub struct {
	gopi.Unit
	gopi.ConnPool
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *stub) New(cfg gopi.Config) error {

	// Return success
	return nil
}
