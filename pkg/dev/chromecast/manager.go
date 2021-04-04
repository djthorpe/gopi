package chromecast

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	// Modules
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
	gopi.Promises

	cast map[string]*Cast
	conn map[string]*Conn
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	serviceTypeCast       = "_googlecast._tcp."
	serviceConnectTimeout = time.Second * 15
	serciceMessageTimeout = time.Second
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	this.Require(this.ServiceDiscovery, this.Logger, this.Publisher)

	// Make map of devices and connections
	this.cast = make(map[string]*Cast)
	this.conn = make(map[string]*Conn)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	// Disconnect devices
	var result error
	for _, cast := range this.cast {
		if err := this.disconnect(cast); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.cast = nil
	this.conn = nil

	// Return any errors
	return result
}

func (this *Manager) Run(ctx context.Context) error {
	// Receive DNS messages for changes in cast status
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	// Loop handling messages until done
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case evt := <-ch:
			if state, ok := evt.(*State); ok && this.Logger.IsDebug() {
				if payload := state.Payload(); payload != nil {
					obj := make(map[string]interface{})
					if err := json.Unmarshal(payload, &obj); err != nil {
						this.Print(err)
					} else if data, err := json.MarshalIndent(obj, "", "  "); err != nil {
						this.Print(err)
					} else {
						this.Print("State: ", state.key, ": ", string(data))
					}
				}
			} else if record, ok := evt.(gopi.ServiceRecord); ok {
				if record.Service() == serviceTypeCast {
					if cast := NewCastFromRecord(record); cast != nil {
						if flags := this.castevent(cast); flags != gopi.CAST_FLAG_NONE {
							if err := this.Publisher.Emit(NewCastEvent(cast, flags), false); err != nil {
								this.Print("CastManager:", err)
							}
						}
					}
				}
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<cast.manager"
	for _, cast := range this.cast {
		str += fmt.Sprint(" ", cast)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) Devices(ctx context.Context) ([]gopi.Cast, error) {
	// Perform the lookup
	records, err := this.ServiceDiscovery.Lookup(ctx, serviceTypeCast)
	if err != nil {
		return nil, err
	}

	// Return any casts found
	result := make([]gopi.Cast, 0, len(records))
	for _, record := range records {
		cast := NewCastFromRecord(record)
		if cast == nil {
			continue
		}

		// Add cast, emit event
		if existing := this.getCastForId(cast.id); existing == nil {
			this.castevent(cast)
		}

		// Append cast onto results
		result = append(result, this.getCastForId(cast.id))
	}

	// Return success
	return result, nil
}

func (this *Manager) Get(key string) gopi.Cast {
	// Get cast by id
	if cast := this.getCastForId(key); cast != nil {
		return cast
	}

	// Iterate through casts
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	key = strings.ToLower(key)
	for _, cast := range this.cast {
		if strings.ToLower(cast.Name()) == key {
			return cast
		}
	}

	// Not found
	return nil
}

func (this *Manager) Connect(ctx context.Context, cast gopi.Cast) error {
	// Check for bad parameters
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("Connect")
	}

	// Device should have been discovered, return outoforder if already connected
	var result error
	if cast := this.getCastForId(cast.Id()); cast == nil {
		return gopi.ErrNotFound.WithPrefix("Connect")
	} else if conn := this.getConnForId(cast.id); conn != nil {
		return gopi.ErrOutOfOrder.WithPrefix("Connect")
	} else if conn, err := cast.ConnectWithTimeout(this.Publisher, serviceConnectTimeout); err != nil {
		return err
	} else {
		this.setConnForId(cast.id, conn)

		// Emit connect message
		if err := this.Publisher.Emit(NewCastEvent(cast, gopi.CAST_FLAG_CONNECT), false); err != nil {
			result = multierror.Append(result, err)
		}

		// Get cast status
		timeout, cancel := context.WithTimeout(ctx, serciceMessageTimeout)
		defer cancel()
		if err := this.Do(timeout, reqGetStatus, conn).Then(this.wait).Finally(this.done, true); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

func (this *Manager) Disconnect(cast gopi.Cast) error {
	// Check for bad parameters
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("Disconnect")
	}

	// Device should have been discovered
	var result error
	if cast := this.getCastForId(cast.Id()); cast == nil {
		return gopi.ErrNotFound.WithPrefix("Disconnect")
	} else if conn := this.getConnForId(cast.id); conn == nil {
		return gopi.ErrOutOfOrder.WithPrefix("Disconnect")
	} else {
		// Send CLOSE message
		if _, data, err := conn.Disconnect(); err != nil {
			result = multierror.Append(result, err)
		} else if err := conn.send(data); err != nil {
			result = multierror.Append(result, err)
		}

		// Close connection
		this.setConnForId(cast.id, nil)
		if err := cast.Disconnect(conn); err != nil {
			result = multierror.Append(result, err)
		}

		// Emit disconnect message
		if err := this.Publisher.Emit(NewCastEvent(cast, gopi.CAST_FLAG_DISCONNECT), false); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

func (this *Manager) SetVolume(ctx context.Context, cast gopi.Cast, level float32) error {
	// Check arguments
	if level < 0.0 {
		level = 0.0
	} else if level > 1.0 {
		level = 1.0
	}

	// If no connection, then connect
	if conn := this.getConnForId(cast.Id()); conn == nil {
		if err := this.Connect(ctx, cast); err != nil {
			return err
		}
	}

	// Send request
	if conn := this.getConnForId(cast.Id()); conn == nil {
		return gopi.ErrInternalAppError.WithPrefix("SetVolume")
	} else {
		timeout, cancel := context.WithTimeout(ctx, serciceMessageTimeout)
		defer cancel()
		return this.Do(timeout, reqSetVolume, []interface{}{conn, level}).Then(this.wait).Finally(this.done, true)
	}
}

func (this *Manager) SetMuted(ctx context.Context, cast gopi.Cast, muted bool) error {
	// If no connection, then connect
	if conn := this.getConnForId(cast.Id()); conn == nil {
		if err := this.Connect(ctx, cast); err != nil {
			return err
		}
	}

	// Send request
	if conn := this.getConnForId(cast.Id()); conn == nil {
		return gopi.ErrInternalAppError.WithPrefix("SetMuted")
	} else {
		timeout, cancel := context.WithTimeout(ctx, serciceMessageTimeout)
		defer cancel()
		return this.Do(timeout, reqSetMuted, []interface{}{conn, muted}).Then(this.wait).Finally(this.done, true)
	}
}

// LaunchAppWithId launches application with Id on a cast device.
func (this *Manager) LaunchAppWithId(ctx context.Context, cast gopi.Cast, app string) error {
	// If no connection, then connect
	if conn := this.getConnForId(cast.Id()); conn == nil {
		if err := this.Connect(ctx, cast); err != nil {
			return err
		}
	}

	// Send request
	if conn := this.getConnForId(cast.Id()); conn == nil {
		return gopi.ErrInternalAppError.WithPrefix("LaunchAppWithId")
	} else {
		timeout, cancel := context.WithTimeout(ctx, serciceMessageTimeout)
		defer cancel()
		return this.Do(timeout, reqLaunchAppWithId, []interface{}{conn, app}).Then(this.wait).Finally(this.done, true)
	}
}

// ConnectMedia initiates a media session
func (this *Manager) ConnectMedia(ctx context.Context, cast gopi.Cast) error {
	// If no connection, then connect
	conn := this.getConnForId(cast.Id())
	if conn == nil {
		if err := this.Connect(ctx, cast); err != nil {
			return err
		} else {
			conn = this.getConnForId(cast.Id())
		}
	}

	if cast := this.getCastForId(cast.Id()); cast == nil {
		return gopi.ErrInternalAppError.WithPrefix("ConnectMedia")
	} else if app := cast.App(); app == nil {
		return gopi.ErrOutOfOrder.WithPrefix("ConnectMedia")
	} else if app.TransportId == "" {
		return gopi.ErrInternalAppError.WithPrefix("ConnectMedia")
	} else if app.IsIdleScreen {
		return gopi.ErrNotFound.WithPrefix("ConnectMedia")
	} else {
		timeout, cancel := context.WithTimeout(ctx, serciceMessageTimeout)
		defer cancel()
		return this.Do(timeout, reqConnectMedia, []interface{}{conn, app.TransportId}).
			Then(reqGetMediaStatus).Then(this.wait).Finally(this.done, true)
	}
}

// DisconnectMedia ends a media session
func (this *Manager) DisconnectMedia(ctx context.Context, cast gopi.Cast) error {
	// If no connection, then connect
	conn := this.getConnForId(cast.Id())
	if conn == nil {
		if err := this.Connect(ctx, cast); err != nil {
			return err
		} else {
			conn = this.getConnForId(cast.Id())
		}
	}

	if cast := this.getCastForId(cast.Id()); cast == nil {
		return gopi.ErrInternalAppError.WithPrefix("DisconnectMedia")
	} else if app := cast.App(); app == nil {
		return gopi.ErrOutOfOrder.WithPrefix("DisconnectMedia")
	} else if app.TransportId == "" {
		return gopi.ErrOutOfOrder.WithPrefix("DisconnectMedia")
	} else if req, data, err := conn.DisconnectMedia(app.TransportId); err != nil {
		return err
	} else if err := conn.send(data); err != nil {
		return err
	} else if flags := cast.UpdateState(NewMediaState(cast.Id(), req, nil, Media{})); flags != gopi.CAST_FLAG_NONE {
		// Emit any media changes
		if err := this.Publisher.Emit(NewCastEvent(cast, flags), false); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

// LoadMedia asks Chromecast to play media
func (this *Manager) LoadMedia(ctx context.Context, cast gopi.Cast, url *url.URL, autoplay bool) error {
	// Check for supported URL schemes
	if url == nil {
		return gopi.ErrBadParameter.WithPrefix("LoadURL")
	} else if url.Scheme != "http" && url.Scheme != "https" {
		return gopi.ErrBadParameter.WithPrefix("Unsupported URL scheme")
	}

	// If no connection, then connect
	if conn := this.getConnForId(cast.Id()); conn == nil {
		if err := this.Connect(ctx, cast); err != nil {
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

	// Get transportId
	var transportId string
	if cast := this.getCastForId(cast.Id()); cast == nil {
		return gopi.ErrInternalAppError.WithPrefix("LoadMedia")
	} else if app := cast.App(); app == nil {
		return gopi.ErrOutOfOrder.WithPrefix("LoadMedia")
	} else if app.TransportId == "" {
		return gopi.ErrOutOfOrder.WithPrefix("LoadMedia")
	} else {
		transportId = app.TransportId
	}

	// Get connection
	conn := this.getConnForId(cast.Id())
	if conn == nil {
		return gopi.ErrInternalAppError.WithPrefix("LaunchAppWithId")
	}

	// Send request
	timeout, cancel := context.WithTimeout(ctx, serciceMessageTimeout)
	defer cancel()
	return this.Do(timeout, reqLoadMedia, []interface{}{conn, transportId, url, mimetype, autoplay}).
		Then(this.wait).Finally(this.done, true)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) disconnect(cast *Cast) error {
	// Check parameters
	if cast == nil {
		return gopi.ErrBadParameter.WithPrefix("Disconnect")
	}

	// Remove cast from list
	if existing := this.getCastForId(cast.id); existing == nil {
		return gopi.ErrNotFound.WithPrefix("Disconnect")
	} else {
		this.setCastForId(cast.id, nil)
	}

	// Remove connection from list
	var result error
	if conn := this.getConnForId(cast.id); conn != nil {
		this.setConnForId(cast.id, nil)
		if err := conn.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

func (this *Manager) getCasts() []*Cast {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	casts := make([]*Cast, 0, len(this.cast))
	for _, cast := range this.cast {
		casts = append(casts, cast)
	}
	return casts
}

func (this *Manager) getCastForId(id string) *Cast {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if cast, exists := this.cast[id]; exists {
		return cast
	} else {
		return nil
	}
}

func (this *Manager) setCastForId(id string, cast *Cast) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if cast == nil {
		delete(this.cast, id)
	} else {
		this.cast[id] = cast
	}
}

func (this *Manager) getConnForId(id string) *Conn {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if conn, exists := this.conn[id]; exists {
		return conn
	} else {
		return nil
	}
}

func (this *Manager) setConnForId(id string, conn *Conn) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if conn == nil {
		delete(this.conn, id)
	} else {
		this.conn[id] = conn
	}
}

// CastEvent returns any changes to a chromecast if it is already
// discovered or returns DISCOVERY flag
func (this *Manager) castevent(cast *Cast) gopi.CastFlag {
	if other := this.getCastForId(cast.id); other == nil {
		this.setCastForId(cast.id, cast)
		return gopi.CAST_FLAG_DISCOVERY
	} else if flags := other.Equals(cast); flags == gopi.CAST_FLAG_NONE {
		return flags
	} else {
		other.updateFrom(cast)
		return flags
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - REQUEST/RESPONSE TO CHROMECAST

type promise struct {
	key string
	req int
}

func reqGetStatus(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("reqGetStatus")
	conn := v.(*Conn)
	if req, data, err := conn.GetStatus(); err != nil {
		return nil, err
	} else if err := conn.send(data); err != nil {
		return nil, err
	} else {
		return &promise{conn.key, req}, nil
	}
}

func reqSetVolume(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("reqSetVolume")
	params := v.([]interface{})
	conn := params[0].(*Conn)
	level := params[1].(float32)
	if req, data, err := conn.SetVolume(Volume{Level: level}); err != nil {
		return nil, err
	} else if err := conn.send(data); err != nil {
		return nil, err
	} else {
		return &promise{conn.key, req}, nil
	}
}

func reqSetMuted(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("reqSetMuted")
	params := v.([]interface{})
	conn := params[0].(*Conn)
	muted := params[1].(bool)
	if req, data, err := conn.SetMuted(muted); err != nil {
		return nil, err
	} else if err := conn.send(data); err != nil {
		return nil, err
	} else {
		return &promise{conn.key, req}, nil
	}
}

func reqLaunchAppWithId(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("reqLaunchAppWithId")
	params := v.([]interface{})
	conn := params[0].(*Conn)
	app := params[1].(string)
	if req, data, err := conn.LaunchAppWithId(app); err != nil {
		return nil, err
	} else if err := conn.send(data); err != nil {
		return nil, err
	} else {
		return &promise{conn.key, req}, nil
	}
}

func reqConnectMedia(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("reqConnectMedia")
	params := v.([]interface{})
	conn := params[0].(*Conn)
	transportId := params[1].(string)
	if _, data, err := conn.ConnectMedia(transportId); err != nil {
		return nil, err
	} else if err := conn.send(data); err != nil {
		return nil, err
	} else {
		return []interface{}{conn, transportId}, nil
	}
}

func reqGetMediaStatus(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("reqGetMediaStatus")
	params := v.([]interface{})
	conn := params[0].(*Conn)
	transportId := params[1].(string)
	if req, data, err := conn.GetMediaStatus(transportId); err != nil {
		return nil, err
	} else if err := conn.send(data); err != nil {
		return nil, err
	} else {
		return &promise{conn.key, req}, nil
	}
}

func reqLoadMedia(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("reqLoadMedia")
	params := v.([]interface{})
	conn := params[0].(*Conn)
	transportId := params[1].(string)
	url := params[2].(*url.URL)
	mimetype := params[3].(string)
	autoplay := params[4].(bool)

	if req, data, err := conn.LoadMedia(transportId, url.String(), mimetype, autoplay); err != nil {
		return nil, err
	} else if err := conn.send(data); err != nil {
		return nil, err
	} else {
		return &promise{conn.key, req}, nil
	}
}

func reqSetPlay(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("reqSetPlay")
	params := v.([]interface{})
	conn := params[0].(*Conn)
	transportId := params[1].(string)
	sessionId := params[2].(int)
	state := params[3].(bool)

	if req, data, err := conn.Play(transportId, sessionId, state); err != nil {
		return nil, err
	} else if err := conn.send(data); err != nil {
		return nil, err
	} else {
		return &promise{conn.key, req}, nil
	}
}

func reqSetPause(ctx context.Context, v interface{}) (interface{}, error) {
	fmt.Println("reqSetPause")
	params := v.([]interface{})
	conn := params[0].(*Conn)
	transportId := params[1].(string)
	sessionId := params[2].(int)
	state := params[3].(bool)

	if req, data, err := conn.Pause(transportId, sessionId, state); err != nil {
		return nil, err
	} else if err := conn.send(data); err != nil {
		return nil, err
	} else {
		return &promise{conn.key, req}, nil
	}
}

func (this *Manager) wait(ctx context.Context, v interface{}) (interface{}, error) {
	// Wait for a response from the chromecast
	promise := v.(*promise)
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case evt := <-ch:
			if evt, ok := evt.(*State); ok {
				if evt.key == promise.key && evt.req == promise.req {
					return evt, nil
				}
			}
		}
	}
}

func (this *Manager) done(v interface{}, err error) error {
	// Process any errors
	if err != nil {
		return err
	}

	// Update state of chromecast or return error from chromecast
	state := v.(*State)
	if state.Err() != nil {
		return state.Err()
	} else if cast := this.getCastForId(state.key); cast != nil {
		if flags := cast.UpdateState(state); flags != gopi.CAST_FLAG_NONE {
			if err := this.Publisher.Emit(NewCastEvent(cast, flags), false); err != nil {
				return err
			}
		}
	}

	// Return success
	return nil
}
