package tradfri

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"path/filepath"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	coap "github.com/go-ocf/go-coap"
	coapnet "github.com/go-ocf/go-coap/net"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	sync.RWMutex
	gopi.Logger
	gopi.Publisher
	Token
	*coap.ClientConn

	// Flags
	key, path *string
	timeout   *time.Duration

	// Path to the configuration
	Path string
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEFAULT_TIMEOUT = 5 * time.Second
	DEFAULT_PORT    = 5684
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) Define(cfg gopi.Config) error {
	this.key = cfg.FlagString("tradfri.key", "", "Tradfri Gateway Key")
	this.timeout = cfg.FlagDuration("tradfri.timeout", DEFAULT_TIMEOUT, "Connection Timeout")
	this.path = cfg.FlagString("tradfri.path", "", "Path to configuration")
	return nil
}

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Publisher, this.Logger)

	// Create config path if it doesn't exist, and read token file
	if path, err := this.Token.CreatePath(*this.path); err != nil {
		return err
	} else if err := this.Token.Read(path); err != nil {
		return err
	} else {
		this.Path = path
	}

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Close connection
	var result error
	if this.ClientConn != nil {
		if err := this.ClientConn.Close(); err != nil && errors.Is(err, coapnet.ErrServerClosed) {
			result = multierror.Append(result, err)
		}
	}

	// Dispose of token
	this.Token.Dispose()

	// Release resources
	this.ClientConn = nil

	// Return any errors
	return result
}

func (this *Manager) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Manager) Addr() net.Addr {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ClientConn == nil {
		return nil
	} else {
		return this.ClientConn.RemoteAddr()
	}
}

func (this *Manager) Id() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ClientConn != nil {
		return this.Token.Id
	} else {
		return ""
	}
}

func (this *Manager) Version() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ClientConn != nil {
		return this.Token.Version
	} else {
		return ""
	}
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT AND DISCONNECT

func (this *Manager) Connect(id string, host string, port uint16) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check state and incoming parameters
	if this.ClientConn != nil {
		return gopi.ErrOutOfOrder.WithPrefix("Connect")
	} else if host == "" || id == "" {
		return gopi.ErrBadParameter.WithPrefix("Connect")
	}

	// Set port if necessary
	if port == 0 {
		port = DEFAULT_PORT
	}

	// Authenticate and Close connection
	addr := fmt.Sprint(host, ":", port)
	if this.Token.Id == "" || this.Token.Id != id || this.Token.Token == "" {
		// Connect
		if conn, err := coapConnectWith(addr, "Client_identity", *this.key, *this.timeout); err != nil {
			return fmt.Errorf("CoapConnect: %w", err)
		} else if response, err := coapAuthenticate(conn, id, *this.timeout); err != nil {
			return fmt.Errorf("CoapAuthenticate: %w", err)
		} else if err := json.Unmarshal(response, &this.Token); err != nil {
			return err
		} else {
			this.Token.Id = id
			if err := this.Token.Write(this.Path); err != nil {
				return err
			} else if err := conn.Close(); err != nil && errors.Is(err, coapnet.ErrServerClosed) {
				return fmt.Errorf("CoapDisconnect: %w", err)
			}
		}
	}

	// Connect with existing token parameters
	if conn, err := coapConnectWith(addr, this.Token.Id, this.Token.Token, *this.timeout); err != nil {
		return err
	} else {
		this.ClientConn = conn
	}

	// Success
	return nil
}

func (this *Manager) Disconnect() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Close connection
	if this.ClientConn != nil {
		if err := this.ClientConn.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.ClientConn = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS DEVICES, GROUPS AND SCENES

func (this *Manager) Devices(ctx context.Context) ([]gopi.TradfriDevice, error) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Check state
	if this.ClientConn == nil {
		return nil, gopi.ErrOutOfOrder.WithPrefix("Devices")
	}

	// Get object id's
	devices, err := coapRequestIdsForPath(ctx, this.ClientConn, ROOT_DEVICES)
	if err != nil {
		return nil, err
	}

	// Get objects
	result := make([]gopi.TradfriDevice, 0, len(devices))
	for _, id := range devices {
		device := new(device)
		if err := coapRequestObjForPath(ctx, this.ClientConn, device, ROOT_DEVICES, fmt.Sprint(id)); err != nil {
			return nil, err
		} else {
			result = append(result, device)
		}
	}

	// Return success
	return result, nil
}

func (this *Manager) ObserveDevice(ctx context.Context, device gopi.TradfriDevice) error {
	var obs *coap.Observation

	// Check state
	if device == nil {
		return gopi.ErrBadParameter.WithPrefix("ObserveDevice")
	} else if this.ClientConn == nil {
		return gopi.ErrOutOfOrder.WithPrefix("ObserveDevice")
	}

	// Observe until cancel
	ticker := time.NewTimer(100 * time.Millisecond)
	path := filepath.Join(append([]string{"/"}, ROOT_DEVICES, fmt.Sprint(device.Id()))...)
	for {
		select {
		case <-ticker.C:
			if obs != nil {
				if err := obs.Cancel(); err != nil {
					return err
				}
			}
			// Start observation
			if obs_, err := this.ClientConn.Observe(path, this.observeDeviceCallback); err != nil {
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
				return nil
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<tradfri.manager"
	if addr := this.Addr(); addr != nil {
		str += fmt.Sprintf(" addr=%q", addr.String())
	}
	str += fmt.Sprintf(" path=%q", this.Path)
	if this.Token.Id != "" {
		str += fmt.Sprint(" token=", this.Token.String())
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) observeDeviceCallback(response *coap.Request) {
	fmt.Println("Callback", response)
}
