package codec

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Codec interface {
	Run(ctx context.Context, publisher gopi.Publisher) error
}

type CodecEvent struct {
	Type   gopi.InputType
	Device gopi.InputDeviceType
	Code   uint32
}

/////////////////////////////////////////////////////////////////////
// PUBLIC PROPERTIES

func (this *CodecEvent) Name() string {
	name := fmt.Sprint(this.Device)
	return strings.ToLower(strings.TrimPrefix(name, "INPUT_DEVICE_"))
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *CodecEvent) String() string {
	str := "<event.codec"
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if this.Type != gopi.INPUT_EVENT_NONE {
		str += " type=" + fmt.Sprint(this.Type)
	}
	if this.Device != gopi.INPUT_DEVICE_NONE {
		str += " device=" + fmt.Sprint(this.Device|gopi.INPUT_DEVICE_REMOTE)
	}
	if this.Code > 0x0000 && this.Code <= 0xFFFF {
		str += fmt.Sprintf(" code=0x%04X", this.Code)
	} else if this.Code > 0x0000 {
		str += fmt.Sprintf(" code=0x%08X", this.Code)
	}
	return str + ">"
}
