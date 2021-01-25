package googlecast

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"sync"
	"syscall"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	sync.RWMutex
	gopi.Unit
	gopi.ServiceDiscovery
	gopi.Publisher
	gopi.Logger

	// Connected Cast Devices
	dev map[string]*Cast

	// Channels for communication
	state chan state
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	serviceTypeCast       = "_googlecast._tcp"
	serviceConnectTimeout = time.Second * 15
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	if this.ServiceDiscovery == nil {
		return gopi.ErrInternalAppError.WithPrefix("ServiceDiscovery")
	}

	// Make map of devices and error channel
	this.dev = make(map[string]*Cast)
	this.state = make(chan state)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Disconnect devices
	var result error
	for _, cast := range this.dev {
		if err := this.disconnect(cast); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Close channels
	close(this.state)

	// Release resources
	this.dev = nil
	this.state = nil

	// Return any errors
	return result
}

func (this *Manager) Run(ctx context.Context) error {
	// Update cast status every second
	timer := time.NewTicker(time.Second)
	defer timer.Stop()

	// Loop handling messages until done
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case state := <-this.state:
			if err := this.setState(state); err != nil {
				this.Print("CastManager: SetState: ", err)
			}
		case <-timer.C:
			if err := this.updateStatus(); err != nil {
				this.Print("CastManager: UpdateStatus: ", err)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Manager) Devices(ctx context.Context) ([]gopi.Cast, error) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Perform the lookup
	records, err := this.ServiceDiscovery.Lookup(ctx, serviceTypeCast)
	if err != nil {
		return nil, err
	}

	result := make([]gopi.Cast, 0, len(records))
	for _, record := range records {
		if cast := NewCastFromRecord(record); cast == nil {
			continue
		} else if connected, exists := this.dev[cast.id]; exists {
			result = append(result, connected)
		} else {
			result = append(result, cast)
		}
	}

	// Return success
	return result, nil
}

func (this *Manager) Connect(device gopi.Cast) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check for bad parameters
	if device == nil {
		return gopi.ErrBadParameter.WithPrefix("Connect")
	}

	// Check for already connected
	key := device.Id()
	if _, exists := this.dev[key]; exists {
		return gopi.ErrDuplicateEntry.WithPrefix("Connect")
	}

	// Do the connection
	if device_, ok := device.(*Cast); ok == false {
		return gopi.ErrInternalAppError.WithPrefix("Connect")
	} else if err := this.connect(device_); err != nil {
		return err
	} else {
		this.dev[key] = device_
	}

	// Emit connect
	if this.Publisher != nil {
		cast := this.dev[key]
		this.Publisher.Emit(NewEvent(cast, nil, nil, gopi.CAST_FLAG_CONNECT, 0), true)
	}

	// Return success
	return nil
}

func (this *Manager) Disconnect(device gopi.Cast) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error
	if device == nil {
		return gopi.ErrBadParameter.WithPrefix("Disconnect")
	}

	key := device.Id()
	if connected, exists := this.dev[key]; exists == false {
		return gopi.ErrNotFound.WithPrefix("Disconnect")
	} else if err := this.disconnect(connected); err != nil {
		result = multierror.Append(result, err)
	}

	// Emit disconnect
	if this.Publisher != nil {
		this.Publisher.Emit(NewEvent(this.dev[key], nil, nil, gopi.CAST_FLAG_DISCONNECT, 0), true)
	}

	// Remove device from list of devices
	delete(this.dev, key)

	// Return any errors
	return result
}

func (this *Manager) Volume(cast gopi.Cast) (float32, bool, error) {
	if cast == nil {
		return 0, false, gopi.ErrBadParameter.WithPrefix("SetVolume")
	}

	device := this.getConnectedDevice(cast)
	if device == nil {
		if err := this.Connect(cast); err != nil {
			return 0, false, err
		}
	}
	if device == nil {
		device = this.getConnectedDevice(cast)
	}

	if device == nil || device.volume == nil {
		return 0, false, nil
	} else {
		return device.volume.Level, device.volume.Muted, nil
	}
}

func (this *Manager) App(cast gopi.Cast) (gopi.CastApp, error) {
	if cast == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("SetVolume")
	}

	device := this.getConnectedDevice(cast)
	if device == nil {
		if err := this.Connect(cast); err != nil {
			return nil, err
		}
	}
	if device == nil {
		device = this.getConnectedDevice(cast)
	}

	if device == nil || device.app == nil || device.app.AppId == "" {
		return nil, nil
	} else {
		return device.app, nil
	}
}

func (this *Manager) LaunchAppWithId(cast gopi.Cast, appId string) error {
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("LaunchAppWithId")
	}
	if appId == "" {
		return gopi.ErrBadParameter.WithPrefix("LaunchAppWithId")
	}

	if device := this.getConnectedDevice(cast); device == nil {
		if err := this.Connect(cast); err != nil {
			return err
		}
	}

	if device := this.getConnectedDevice(cast); device == nil {
		return gopi.ErrNotFound.WithPrefix("LaunchAppWithId")
	} else {
		return device.ReqLaunchAppWithId(appId)
	}
}

func (this *Manager) SetVolume(cast gopi.Cast, value float32) error {
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("SetVolume")
	}

	if device := this.getConnectedDevice(cast); device == nil {
		if err := this.Connect(cast); err != nil {
			return err
		}
	}

	if device := this.getConnectedDevice(cast); device == nil {
		return gopi.ErrNotFound.WithPrefix("SetVolume")
	} else {
		return device.ReqVolumeLevel(value)
	}
}

func (this *Manager) SetMuted(cast gopi.Cast, value bool) error {
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("SetMuted")
	}

	if device := this.getConnectedDevice(cast); device == nil {
		if err := this.Connect(cast); err != nil {
			return err
		}
	}

	if device := this.getConnectedDevice(cast); device == nil {
		return gopi.ErrNotFound.WithPrefix("SetMuted")
	} else {
		return device.ReqMuted(value)
	}
}

func (this *Manager) SetPlay(cast gopi.Cast, value bool) error {
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("SetPlay")
	}

	if device := this.getConnectedDevice(cast); device == nil {
		if err := this.Connect(cast); err != nil {
			return err
		}
	}

	if device := this.getConnectedDevice(cast); device == nil {
		return gopi.ErrNotFound.WithPrefix("SetPlay")
	} else {
		return device.ReqPlay(value)
	}
}

func (this *Manager) SetPause(cast gopi.Cast, value bool) error {
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("SetPause")
	}

	if device := this.getConnectedDevice(cast); device == nil {
		if err := this.Connect(cast); err != nil {
			return err
		}
	}

	if device := this.getConnectedDevice(cast); device == nil {
		return gopi.ErrNotFound.WithPrefix("SetPause")
	} else {
		return device.ReqPause(value)
	}
}

func (this *Manager) SeekAbs(cast gopi.Cast, value time.Duration) error {
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("SeekAbs")
	}

	if device := this.getConnectedDevice(cast); device == nil {
		if err := this.Connect(cast); err != nil {
			return err
		}
	}

	if device := this.getConnectedDevice(cast); device == nil {
		return gopi.ErrNotFound.WithPrefix("SeekAbs")
	} else {
		return device.ReqSeekAbs(float32(value.Seconds()))
	}
}

func (this *Manager) SeekRel(cast gopi.Cast, value time.Duration) error {
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("SeekRel")
	}

	if device := this.getConnectedDevice(cast); device == nil {
		if err := this.Connect(cast); err != nil {
			return err
		}
	}

	if device := this.getConnectedDevice(cast); device == nil {
		return gopi.ErrNotFound.WithPrefix("SeekRel")
	} else {
		return device.ReqSeekRel(float32(value.Seconds()))
	}
}

func (this *Manager) Stop(cast gopi.Cast) error {
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("Stop")
	}

	if device := this.getConnectedDevice(cast); device == nil {
		if err := this.Connect(cast); err != nil {
			return err
		}
	}

	if device := this.getConnectedDevice(cast); device == nil {
		return gopi.ErrNotFound.WithPrefix("Stop")
	} else {
		return device.ReqStop()
	}
}

func (this *Manager) LoadURL(cast gopi.Cast, url *url.URL, autoplay bool) error {
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("LoadURL")
	}
	if url == nil {
		return gopi.ErrBadParameter.WithPrefix("LoadURL")
	} else if url.Scheme != "http" && url.Scheme != "https" {
		return gopi.ErrBadParameter.WithPrefix("Unsupported URL scheme")
	}

	// Connect device
	if device := this.getConnectedDevice(cast); device == nil {
		if err := this.Connect(cast); err != nil {
			return err
		}
	}

	// Get mimetype
	mimetype := "application/octet-stream"
	skipverify := true
	client := http.Client{
		Timeout: serviceConnectTimeout,
	}
	client.Transport = http.DefaultTransport
	if skipverify {
		client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	if response, err := client.Head(url.String()); err != nil {
		return err
	} else if response.StatusCode != http.StatusOK {
		return gopi.ErrUnexpectedResponse.WithPrefix(response.Status)
	} else if mimetype_ := response.Header.Get("Content-Type"); mimetype == "" {
		return gopi.ErrUnexpectedResponse.WithPrefix("Content-Type")
	} else if contenttype, _, err := mime.ParseMediaType(mimetype_); err != nil {
		return err
	} else {
		mimetype = contenttype
	}

	// Request load
	if device := this.getConnectedDevice(cast); device == nil {
		return gopi.ErrNotFound.WithPrefix("LoadURL")
	} else {
		return device.ReqLoadURL(url, mimetype, autoplay)
	}

}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<cast.manager"
	for _, device := range this.dev {
		str += fmt.Sprint(" %v=%v", device.Id(), device)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) disconnect(device *Cast) error {
	return device.Disconnect()
}

func (this *Manager) connect(device *Cast) error {
	return device.ConnectWithTimeout(serviceConnectTimeout, this.state)
}

func (this *Manager) getConnectedDevice(cast gopi.Cast) *Cast {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	key := cast.Id()
	if dev, exists := this.dev[key]; exists {
		return dev
	} else {
		return nil
	}
}

func (this *Manager) updateStatus() error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	var result error
	for _, device := range this.dev {
		if device.isConnected() == false {
			// ignore device
		} else if err := device.UpdateState(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

func (this *Manager) setState(s state) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Find device to change state on
	device, exists := this.dev[s.key]
	if exists == false {
		this.Debug("Debug: ", device.Id(), ": ", gopi.ErrNotFound)
		return nil
	}

	// Check for error
	if s.err != nil {
		this.Print("Error: ", device.Id(), ": ", s.err)
		return nil
	}

	// Output any debug messages
	if s.dbg != "" {
		this.Debug("Debug: ", s.key, ": ", s.dbg)
	}

	// If close, then disconnect
	if s.close {
		this.Print(device.Id(), ": Disconnecting after CLOSE message received")
		go this.Disconnect(device)
	}

	// Return if no state
	if len(s.values) == 0 {
		return nil
	}

	// Set state in device
	if flags, err := device.SetState(s); err != nil {
		return err
	} else if flags != gopi.CAST_FLAG_NONE && this.Publisher != nil {
		this.Publisher.Emit(NewEvent(device, device.app, device.volume, flags, s.req), true)

		if flags&gopi.CAST_FLAG_APP != 0 {
			this.Debug("App:", device.id, "=>", device.app)
		}
		if flags&gopi.CAST_FLAG_VOLUME != 0 {
			this.Debug("Vol:", device.id, "=>", device.volume)
		}
		if flags&gopi.CAST_FLAG_MEDIA != 0 {
			this.Debug("Media:", device.id, "=>", device.media)
		}
	}

	// Return success
	return nil
}

func isDroppedConnection(err error) bool {
	if errors.Is(err, syscall.ECONNABORTED) {
		return true
	} else if errors.Is(err, syscall.ECONNREFUSED) {
		return true
	} else if errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	return false
}
