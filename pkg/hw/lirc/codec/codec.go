package codec

import (
	"context"

	"github.com/djthorpe/gopi/v3"
)

type Codec interface {
	Run(ctx context.Context, publisher gopi.Publisher) error
}

type CodecEvent struct {
	Type     gopi.InputType
	Device   gopi.InputDevice
	Scancode uint32
}

func (this *CodecEvent) Name() string {
	return "CodecEvent"
}
