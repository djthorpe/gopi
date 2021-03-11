package rotel

import (
	"context"
	"io"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
	"github.com/golang/protobuf/ptypes/empty"

	//empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type stub struct {
	gopi.Conn
	ManagerClient
}

/////////////////////////////////////////////////////////////////////
// INIT

func (this *stub) New(conn gopi.Conn) {
	this.Conn = conn
	this.ManagerClient = NewManagerClient(conn.(grpc.ClientConnInterface))
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *stub) SetPower(ctx context.Context, state bool) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.SetPower(ctx, &Bool{Value: state}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) SetSource(ctx context.Context, source string) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.SetSource(ctx, &String{Value: source}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) SetVolume(ctx context.Context, value uint) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.SetVolume(ctx, &Uint{Value: uint32(value)}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) SetMute(ctx context.Context, value bool) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.SetMute(ctx, &Bool{Value: value}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) SetBypass(ctx context.Context, value bool) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.SetBypass(ctx, &Bool{Value: value}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) SetTreble(ctx context.Context, value int) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.SetTreble(ctx, &Int{Value: int32(value)}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) SetBass(ctx context.Context, value int) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.SetBass(ctx, &Int{Value: int32(value)}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) SetBalance(ctx context.Context, location string, value uint) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.SetBalance(ctx, &Balance{Location: location, Value: uint32(value)}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) SetDimmer(ctx context.Context, value uint) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.SetDimmer(ctx, &Uint{Value: uint32(value)}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) Play(ctx context.Context) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.Play(ctx, &empty.Empty{}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) Stop(ctx context.Context) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.Stop(ctx, &empty.Empty{}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) Pause(ctx context.Context) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.Pause(ctx, &empty.Empty{}); err != nil {
		return err
	} else {
		return nil
	}
}
func (this *stub) NextTrack(ctx context.Context) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.NextTrack(ctx, &empty.Empty{}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) PrevTrack(ctx context.Context) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	if _, err := this.ManagerClient.PrevTrack(ctx, &empty.Empty{}); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *stub) Stream(ctx context.Context, ch chan<- gopi.RotelEvent) error {
	this.Conn.Lock()
	defer this.Conn.Unlock()

	stream, err := this.ManagerClient.Stream(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if msg, err := stream.Recv(); err == io.EOF {
				return nil
			} else if err != nil {
				return this.Err(err)
			} else if evt := fromProtoEvent(msg); evt != nil {
				ch <- evt
			}
		}
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *stub) String() string {
	str := "<rpc.rotelstub"
	str += " addr=" + strconv.Quote(this.Addr())
	return str + ">"
}
