/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	coap "github.com/go-ocf/go-coap"
	"github.com/go-ocf/go-coap/codes"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type gateway struct {
	base.Unit
	sync.Mutex
	token

	key     string
	path    string
	addr    string
	timeout time.Duration
	conn    *coap.ClientConn
	devices map[uint]*device
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

const (
	CONN_TIMEOUT       = 5 * time.Second
	PATH_AUTH_EXCHANGE = "/15011/9063"
	PATH_DEVICES       = "/15001"
	PATH_GROUPS        = "/15004"
	PATH_SCENES        = "/15005"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *gateway) Init(config Tradfri) error {
	this.token.Id = config.Id
	this.key = config.Key

	// Set timeout
	if config.Timeout == 0 {
		this.timeout = CONN_TIMEOUT
	} else {
		this.timeout = config.Timeout
	}

	// Create path if it doesn't exist, and read token file
	if path, err := this.token.CreatePath(config.Path); err != nil {
		return err
	} else if err := this.token.Read(path); err != nil {
		return err
	} else {
		this.path = path
	}

	// Create empty devices map
	this.devices = make(map[uint]*device)

	// Success
	return nil
}

func (this *gateway) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Close connection
	//if this.conn != nil {
	//if err := this.conn.Close(); err != nil {
	//	return err
	//}
	//}

	// Release resources
	this.conn = nil
	this.devices = nil

	// Success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *gateway) String() string {
	str := "<tradfri.Gateway " + this.token.String()
	if this.addr != "" {
		str += " addr=" + strconv.Quote(this.addr)
	}
	if this.key != "" {
		str += " key=" + strconv.Quote(this.key)
	}
	if this.path != "" {
		str += " path=" + strconv.Quote(this.path)
	}
	if this.timeout != 0 {
		str += " timeout=" + fmt.Sprint(this.timeout)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT

func (this *gateway) Connect(service gopi.RPCServiceRecord, flags gopi.RPCFlag) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.conn != nil {
		return gopi.ErrOutOfOrder
	} else if addr, err := addrForService(service, flags); err != nil {
		return err
	} else {
		// Authenticate and close connection
		if this.token.Id == "" || this.token.Token == "" {
			if id := strings.TrimSpace(service.Name); id == "" {
				return gopi.ErrBadParameter.WithPrefix("id")
			} else if conn, err := coapConnectWith(addr, "Client_identity", this.key, this.timeout); err != nil {
				return err
			} else if response, err := coapAuthenticate(conn, id, this.timeout); err != nil {
				return err
			} else if err := json.Unmarshal(response, &this.token); err != nil {
				return err
			} else {
				this.token.Id = id
				if err := this.token.Write(this.path); err != nil {
					return err
				} else if err := conn.Close(); err != nil {
					return err
				}
			}
		}

		// Connect with existing token parameters
		if conn, err := coapConnectWith(addr, this.token.Id, this.token.Token, this.timeout); err != nil {
			return err
		} else {
			this.conn = conn
			this.addr = addr
		}
	}

	// Emit connected message
	// TODO this.bus.Emit(NewEvent(this, mutablehome.IKEA_EVENT_GATEWAY_CONNECTED, nil))

	// Success
	return nil
}

func (this *gateway) Disconnect() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Emit disconnected message
	// TODO this.bus.Emit(NewEvent(this, mutablehome.IKEA_EVENT_GATEWAY_DISCONNECTED, nil))

	// Close connection
	if this.conn != nil {
		if err := this.conn.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.conn = nil
	this.addr = ""

	// Return success
	return nil
}

func (this *gateway) Id() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.conn != nil {
		return this.token.Id
	} else {
		return ""
	}
}

func (this *gateway) Version() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.conn != nil {
		return this.token.Version
	} else {
		return ""
	}
}

////////////////////////////////////////////////////////////////////////////////
// DEVICES, GROUPS AND SCENES

func (this *gateway) Devices() ([]uint, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.requestIdsForPath(PATH_DEVICES)
}

func (this *gateway) Groups() ([]uint, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.requestIdsForPath(PATH_GROUPS)
}

func (this *gateway) Scenes() ([]uint, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.requestIdsForPath(PATH_SCENES)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *gateway) requestIdsForPath(path string) ([]uint, error) {
	if this.conn == nil {
		return nil, gopi.ErrOutOfOrder
	}
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()

	var ids []uint
	if response, err := this.conn.GetWithContext(ctx, path); err != nil {
		return nil, err
	} else if response.Code() != codes.Content {
		return nil, fmt.Errorf("%w: %v (path: %v)", gopi.ErrUnexpectedResponse, response.Code(), strconv.Quote(path))
	} else if err := json.Unmarshal(response.Payload(), &ids); err != nil {
		return nil, fmt.Errorf("%w: %v", err, string(response.Payload()))
	}

	// Success
	return ids, nil
}

func (this *gateway) requestObjForPathId(path string, id uint, obj interface{}) error {
	if this.conn == nil {
		return gopi.ErrOutOfOrder
	}
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()
	if response, err := this.conn.GetWithContext(ctx, fmt.Sprintf("%v/%d", path, id)); err != nil {
		return err
	} else if response.Code() != codes.Content {
		return fmt.Errorf("%w: %v (path: %v)", gopi.ErrUnexpectedResponse, response.Code(), response.Path())
	} else if err := json.Unmarshal(response.Payload(), obj); err != nil {
		return fmt.Errorf("%w: %v", err, string(response.Payload()))
	} else {
		fmt.Println(string(response.Payload()))
	}

	// Success
	return nil
}

/*
func (this *tradfri) Device(id uint) (mutablehome.IkeaDevice, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	device := NewDevice()
	if err := this.requestObjForPathId(PATH_DEVICES, id, device); err != nil {
		return nil, err
	} else {
		return device, nil
	}
}

func (this *tradfri) Group(id uint) (mutablehome.IkeaGroup, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	group := NewGroup()
	if err := this.requestObjForPathId(PATH_GROUPS, id, group); err != nil {
		return nil, err
	} else {
		return group, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// SEND COMMANDS TO GATEWAY

func (this *tradfri) Send(commands ...mutablehome.IkeaCommand) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(commands) == 0 {
		return gopi.ErrBadParameter.WithPrefix("values")
	}
	for _, command := range commands {
		if command == nil {
			return gopi.ErrBadParameter.WithPrefix("command")
		} else if body, err := command.Body(); err != nil {
			return gopi.ErrBadParameter.WithPrefix("command")
		} else if message, err := this.conn.Put(command.Path(), coap.AppJSON, body); err != nil {
			return err
		} else if message.Code() != codes.Changed {
			return fmt.Errorf("%w: %v (path: %v)", gopi.ErrUnexpectedResponse, message.Code(), message.Path())
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// OBSERVE DEVICE CHANGES

func (this *tradfri) ObserveDevice(ctx context.Context, id uint) error {
	path := strings.Join([]string{PATH_DEVICES, fmt.Sprint(id)}, "/")
	ticker := time.NewTimer(100 * time.Millisecond)
	var obs *coap.Observation

FOR_LOOP:
	for {
		select {
		case <-ticker.C:
			// Stop
			if obs != nil {
				if err := obs.Cancel(); err != nil {
					return err
				}
			}
			// Start
			if obs_, err := this.conn.Observe(path, this.observeDeviceCallback); err != nil {
				return err
			} else {
				obs = obs_
			}
			// Restart the ticker with some random additional interval
			ticker.Reset(time.Second * (5 + time.Duration(rand.Int31n(15))))
		case <-ctx.Done():
			if err := obs.Cancel(); err != nil {
				return err
			} else {
				break FOR_LOOP
			}
		}
	}

	// Success
	return ctx.Err()
}

func (this *tradfri) observeDeviceCallback(response *coap.Request) {
	device := NewDevice()
	if response.Msg.Code() != codes.Content {
		this.Log.Error(fmt.Errorf("%w: %v", gopi.ErrUnexpectedResponse, response.Msg.Code()))
		return
	} else if err := json.Unmarshal(response.Msg.Payload(), &device); err != nil {
		this.Log.Error(fmt.Errorf("%w: %v", gopi.ErrUnexpectedResponse, err))
		return
	}

	// Emit device event if there is a change
	if event := this.observeDeviceEvent(device); event != mutablehome.IKEA_EVENT_NONE {
		this.bus.Emit(NewEvent(this, event, device))
	}
}

func (this *tradfri) observeDeviceEvent(device *device) mutablehome.IkeaEventType {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device == nil || device.Id() == 0 {
		return mutablehome.IKEA_EVENT_NONE
	}

	// Check for added devices
	id := device.Id()
	if other, exists := this.devices[id]; exists == false {
		this.devices[id] = device
		return mutablehome.IKEA_EVENT_DEVICE_ADDED
	} else if device.Equals(other) == false {
		this.devices[id] = device
		return mutablehome.IKEA_EVENT_DEVICE_CHANGED
	} else {
		return mutablehome.IKEA_EVENT_NONE
	}
}
*/
