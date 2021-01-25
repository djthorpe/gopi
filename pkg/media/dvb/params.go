// +build dvb

package dvb

import (
	dvb "github.com/djthorpe/gopi/v3/pkg/sys/dvb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Params struct {
	*dvb.TuneParams
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewParams(v *dvb.TuneParams) *Params {
	this := new(Params)

	if v == nil {
		return nil
	} else {
		this.TuneParams = v
	}

	// Return success
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Params) Name() string {
	return this.TuneParams.Name()
}
