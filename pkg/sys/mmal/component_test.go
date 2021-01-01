//+build mmal

package mmal_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/mmal"
)

func Test_Component_001(t *testing.T) {
	if comp, err := mmal.MMALComponentCreate(mmal.MMAL_COMPONENT_DEFAULT_IMAGE_ENCODER); err != nil {
		t.Error(err)
	} else {
		t.Log(comp)
		comp.Free()
	}
}
