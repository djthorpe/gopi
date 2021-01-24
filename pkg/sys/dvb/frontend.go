// +build dvb

package dvb

import (
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: FRONT END

func FEPath(Device, frontend uint) string {
	return fmt.Sprintf("%v%v/frontend%v", DVB_PATH_WILDCARD, bus, frontend)
}

func DVB_FEOpen(bus, frontend uint) (*os.File, error) {
	if file, err := os.OpenFile(DVB_FEPath(bus, frontend), os.O_SYNC|os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}
