package googlecast

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Cast struct {
	sync.RWMutex
	connection
	promises

	// Data about the cast device
	id, fn string
	md, rs string
	st     uint
	ips    []net.IP
	port   uint16

	// State information
	volume *Volume
	app    *App
	media  *MediaItem
	player *Media

	// Promises
	callbacks map[int]promise
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	promiseTimeout = 2 * time.Second
	pingTimeout    = 20 * time.Second
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCastFromRecord(r gopi.ServiceRecord) *Cast {
	this := new(Cast)

	// Set addr and port
	if port := r.Port(); port == 0 {
		return nil
	} else {
		this.port = port
	}
	if ips := r.Addrs(); len(ips) == 0 {
		return nil
	} else {
		this.ips = ips
	}

	// Set properties
	tuples := txtToMap(r.Txt())
	if id, exists := tuples["id"]; exists && id != "" {
		this.id = id
	} else {
		return nil
	}
	if fn, exists := tuples["fn"]; exists && fn != "" {
		this.fn = fn
	} else {
		this.fn = this.id
	}
	if md, exists := tuples["md"]; exists {
		this.md = md
	}
	if rs, exists := tuples["rs"]; exists {
		this.rs = rs
	}
	if st, exists := tuples["st"]; exists {
		if st, err := strconv.ParseUint(st, 0, 64); err == nil {
			this.st = uint(st)
		}
	}

	return this
}

func (this *Cast) ConnectWithTimeout(timeout time.Duration, state chan<- state) error {

	// Get an address to connect to
	if len(this.ips) == 0 {
		return gopi.ErrNotFound.WithPrefix("ConnectWithTimeout", "No Address")
	}

	// Use first IP
	// TODO: Use a random IP and retry with other IP's if not working
	addr := fmt.Sprintf("%v:%v", this.ips[0], this.port)
	if err := this.connection.Connect(this.Id(), addr, timeout, state); err != nil {
		return err
	}

	// Lock for setting state
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Set state
	this.volume = nil
	this.app = nil
	this.media = nil
	this.player = nil

	// Init promises
	this.promises.InitWithTimeout(promiseTimeout)

	// Return success
	return nil
}

func (this *Cast) Disconnect() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Disconnect
	if err := this.connection.Disconnect(); err != nil {
		result = multierror.Append(result, err)
	}

	// Set state
	this.volume = nil
	this.app = nil
	this.media = nil
	this.player = nil

	// Dispose promises
	this.promises.Dispose()

	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

// Id returns the identifier for a chromecast
func (this *Cast) Id() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.id
}

// Name returns the readable name for a chromecast
func (this *Cast) Name() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.fn
}

// Model returns the reported model information
func (this *Cast) Model() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.md
}

// Service returns the currently running service
func (this *Cast) Service() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if this.app != nil && this.app.DisplayName != "" {
		return this.app.DisplayName
	} else {
		return this.rs
	}
}

// State returns 0 if backdrop, else returns 1
func (this *Cast) State() uint {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.st
}

////////////////////////////////////////////////////////////////////////////////
// STATE

func (this *Cast) UpdateState() error {
	// Update status for volume and app
	this.RWMutex.RLock()
	if this.volume == nil || this.app == nil {
		this.Debugf("Requesting Volume and App State")
		if _, data, err := this.channel.GetStatus(); err != nil {
			return err
		} else if err := this.send(data); err != nil {
			return err
		}
	} else if this.player == nil && this.app.TransportId != "" && this.app.IsIdleScreen == false {
		this.Debugf("Connecting Media")
		if _, data, err := this.channel.ConnectMedia(this.app.TransportId); err != nil {
			return err
		} else if err := this.send(data); err != nil {
			return err
		}

		this.Debugf("Get Media Status")
		if _, data, err := this.channel.GetMediaStatus(this.app.TransportId); err != nil {
			return err
		} else if err := this.send(data); err != nil {
			return err
		}
	}
	this.RWMutex.RUnlock()

	// If no ping/pong has been done recently then disconnect
	if this.channel.ping.IsZero() == false && time.Since(this.channel.ping) > pingTimeout {
		this.Debugf("Stale Ping, Disconnecting")
		go func() {
			this.ch <- Close(this.channel.key)
		}()
	}

	// Return success
	return nil
}

func (this *Cast) SetState(s state) (gopi.CastFlag, error) {
	// Set state on device
	flags := gopi.CAST_FLAG_NONE
	for _, value := range s.values {
		switch value := value.(type) {
		case Volume:
			this.Debugf("SetState: Volume: %v", value)
			if f, err := this.SetVolume(value); err != nil {
				return flags, err
			} else {
				flags |= f
			}
		case App:
			this.Debugf("SetState: App: %v", value)
			if f, err := this.SetApp(value); err != nil {
				return flags, err
			} else {
				flags |= f
			}
			if this.app == nil || this.app.IsIdleScreen {
				if f, err := this.SetMedia(Media{}); err != nil {
					return flags, err
				} else {
					flags |= f
				}
			}
		case Media:
			this.Debugf("SetState: Media: %v", value)
			if f, err := this.SetMedia(value); err != nil {
				return flags, err
			} else {
				flags |= f
			}
		default:
			return flags, gopi.ErrInternalAppError.WithPrefix(value)
		}
	}

	// Call promise for request
	if err := this.promises.Call(s.req); err != nil {
		this.Debugf("Error: %v", err)
	}

	// Return success
	return flags, nil
}

func (this *Cast) SetVolume(v Volume) (gopi.CastFlag, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.volume == nil || this.volume.Equals(v) == false {
		this.volume = &v
		return gopi.CAST_FLAG_VOLUME, nil
	} else {
		return gopi.CAST_FLAG_NONE, nil
	}
}

func (this *Cast) SetApp(a App) (gopi.CastFlag, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.app == nil || this.app.Equals(a) == false {
		this.app = &a
		return gopi.CAST_FLAG_APP, nil
	} else {
		return gopi.CAST_FLAG_NONE, nil
	}
}

func (this *Cast) SetMedia(m Media) (gopi.CastFlag, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.player == nil || this.player.Equals(m) == false {
		this.player = &m
		// TODO: Equals
		return gopi.CAST_FLAG_MEDIA, nil
	} else {
		return gopi.CAST_FLAG_NONE, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - REQUESTS

func (this *Cast) ReqLaunchAppWithId(appId string) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if _, data, err := this.channel.LaunchAppWithId(appId); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Cast) ReqVolumeLevel(level float32) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Clamp value between 0.0 and 1.0
	if level < 0.0 {
		level = 0.0
	} else if level > 1.0 {
		level = 1.0
	}
	v := Volume{level, false}
	if level == 0 {
		v = Volume{0, true}
	}

	if _, data, err := this.channel.SetVolume(v); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Cast) ReqMuted(muted bool) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if _, data, err := this.channel.SetMuted(muted); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Cast) ReqPlay(state bool) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.media == nil || this.player.MediaSessionId == 0 {
		// Connect media and then provide callback
		return this.ReqMediaConnect(func() error {
			return this.ReqPlay(state)
		})
	} else if _, data, err := this.channel.Play(this.app.TransportId, this.player.MediaSessionId, state); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Cast) ReqPause(state bool) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.media == nil || this.player.MediaSessionId == 0 {
		// Connect media and then provide callback
		return this.ReqMediaConnect(func() error {
			if this.media == nil || this.player.MediaSessionId == 0 {
				return gopi.ErrNotFound.WithPrefix("No Media Playing")
			} else {
				return this.ReqPause(state)
			}
		})
	} else if _, data, err := this.channel.Pause(this.app.TransportId, this.player.MediaSessionId, state); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Cast) ReqSeekAbs(value float32) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.media == nil || this.player.MediaSessionId == 0 {
		// Connect media and then provide callback
		return this.ReqMediaConnect(func() error {
			if this.media == nil || this.player.MediaSessionId == 0 {
				return gopi.ErrNotFound.WithPrefix("No Media Playing")
			} else {
				return this.ReqSeekAbs(value)
			}
		})
	} else if _, data, err := this.channel.SeekAbs(this.app.TransportId, this.player.MediaSessionId, value); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Cast) ReqSeekRel(value float32) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.media == nil || this.player.MediaSessionId == 0 {
		// Connect media and then provide callback
		return this.ReqMediaConnect(func() error {
			if this.media == nil || this.player.MediaSessionId == 0 {
				return gopi.ErrNotFound.WithPrefix("No Media Playing")
			} else {
				return this.ReqSeekRel(value)
			}
		})
	} else if _, data, err := this.channel.SeekRel(this.app.TransportId, this.player.MediaSessionId, value); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Cast) ReqStop() error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if _, data, err := this.channel.Stop(); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Cast) ReqLoadURL(url *url.URL, mimetype string, autoplay bool) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// ConnectMedia
	if this.app == nil || this.app.TransportId == "" {
		return this.ReqAppConnect(func() error {
			return this.ReqLoadURL(url, mimetype, autoplay)
		})
	} else if _, data, err := this.channel.ConnectMedia(this.app.TransportId); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// LoadUrl
	if _, data, err := this.channel.LoadUrl(this.app.TransportId, url.String(), mimetype, autoplay); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Cast) ReqAppConnect(fn func() error) error {
	if reqId, data, err := this.channel.GetStatus(); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	} else {
		// Set promise in background
		go this.promises.Set(reqId, fn)
	}
	// Return success
	return nil
}

func (this *Cast) ReqMediaConnect(fn func() error) error {
	if this.app == nil || this.app.TransportId == "" {
		return this.ReqAppConnect(func() error {
			return this.ReqMediaConnect(fn)
		})
	} else if _, data, err := this.channel.ConnectMedia(this.app.TransportId); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	} else if reqId, data, err := this.channel.GetMediaStatus(this.app.TransportId); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	} else {
		// Set promise in background
		go this.promises.Set(reqId, fn)
	}
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SEND DEBUG MESSAGES

func (this *Cast) Debugf(format string, a ...interface{}) {
	fmt.Printf("Cast %q: %v\n", this.id, fmt.Sprintf(format, a...))
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Cast) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<cast.device"
	str += " id=" + this.Id()
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if model := this.Model(); model != "" {
		str += " model=" + strconv.Quote(model)
	}
	if service := this.Service(); service != "" {
		str += " service=" + strconv.Quote(service)
	}
	str += " state=" + fmt.Sprint(this.State())
	if this.volume != nil {
		str += " volume=" + fmt.Sprint(this.volume)
	}
	if this.app != nil {
		str += " app=" + fmt.Sprint(this.app)
	}
	if this.player != nil {
		str += " player=" + fmt.Sprint(this.player)
	}
	if this.media != nil {
		str += " media=" + fmt.Sprint(this.media)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func txtToMap(txt []string) map[string]string {
	result := make(map[string]string, len(txt))
	for _, r := range txt {
		if kv := strings.SplitN(r, "=", 2); len(kv) == 2 {
			result[kv[0]] = kv[1]
		} else if len(kv) == 1 {
			result[kv[0]] = ""
		}
	}
	return result
}
