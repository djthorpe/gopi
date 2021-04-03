package chromecast

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	pb "github.com/djthorpe/gopi/v3/pkg/rpc/castchannel"
	proto "github.com/golang/protobuf/proto"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Channel struct {
	sync.RWMutex
	gopi.Publisher

	msg  int
	key  string
	ping time.Time
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	CAST_DEFAULT_SENDER   = "sender-0"
	CAST_DEFAULT_RECEIVER = "receiver-0"
	CAST_NS_CONN          = "urn:x-cast:com.google.cast.tp.connection"
	CAST_NS_HEARTBEAT     = "urn:x-cast:com.google.cast.tp.heartbeat"
	CAST_NS_RECV          = "urn:x-cast:com.google.cast.receiver"
	CAST_NS_MEDIA         = "urn:x-cast:com.google.cast.media"
	CAST_NS_MULTIZONE     = "urn:x-cast:com.google.cast.multizone"
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Channel) Init(ch gopi.Publisher, key string) {
	this.Publisher = ch
	this.key = key
}

// Generate a connect message
func (this *Channel) Connect() (int, []byte, error) {
	payload := &PayloadHeader{Type: "CONNECT"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_CONN, payload.WithId(id))

	// Reset ping time
	this.ping = time.Time{}

	// Return payload
	return id, data, err
}

// Generate a disconnect message
func (this *Channel) Disconnect() (int, []byte, error) {
	payload := &PayloadHeader{Type: "CLOSE"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

// GetStatus message
func (this *Channel) GetStatus() (int, []byte, error) {
	payload := &PayloadHeader{Type: "GET_STATUS"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

// Volume
func (this *Channel) SetVolume(v Volume) (int, []byte, error) {
	payload := &SetVolumeRequest{PayloadHeader{Type: "SET_VOLUME"}, v}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

// Mute
func (this *Channel) SetMuted(muted bool) (int, []byte, error) {
	payload := &SetVolumeRequest{PayloadHeader{Type: "SET_VOLUME"}, Volume{Muted: muted}}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

// Launch application
func (this *Channel) LaunchAppWithId(appId string) (int, []byte, error) {
	payload := &LaunchAppRequest{PayloadHeader{Type: "LAUNCH"}, appId}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

// Load media
func (this *Channel) LoadMedia(transportId string, url, mimetype string, autoplay bool) (int, []byte, error) {
	payload := &LoadMediaRequest{}
	payload.PayloadHeader = PayloadHeader{Type: "LOAD"}
	payload.Autoplay = autoplay
	payload.Media.ContentId = url
	payload.Media.ContentType = mimetype
	payload.Media.StreamType = "BUFFERED"
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_MEDIA, payload.WithId(id))
	return id, data, err
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

// PingTime returns the how long last ping received from a chromecast
func (this *Channel) PingTime() time.Duration {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if this.ping.IsZero() {
		return 0
	} else {
		return time.Now().Sub(this.ping)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// encode message and return it
func (this *Channel) encode(source, dest, ns string, payload Payload) ([]byte, error) {
	json, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	payloadStr := string(json)
	message := &pb.CastMessage{
		ProtocolVersion: pb.CastMessage_CASTV2_1_0.Enum(),
		SourceId:        &source,
		DestinationId:   &dest,
		Namespace:       &ns,
		PayloadType:     pb.CastMessage_STRING.Enum(),
		PayloadUtf8:     &payloadStr,
	}
	proto.SetDefaults(message)
	return proto.Marshal(message)
}

// decode message and process it, may return a message which
// needs to be sent in return (for heartbeat messages)
func (this *Channel) decode(data []byte) ([]byte, error) {
	message := &pb.CastMessage{}
	if err := proto.Unmarshal(data, message); err != nil {
		return nil, err
	}
	ns := message.GetNamespace()
	switch ns {
	case CAST_NS_RECV:
		return this.rcvReceiver(message)
	case CAST_NS_HEARTBEAT:
		return this.recvHeartbeat(message)
	case CAST_NS_CONN:
		return this.rcvConnection(message)
	case CAST_NS_MEDIA:
		return this.rcvMedia(message)
	case CAST_NS_MULTIZONE:
		// Ignore messages for Multizone, I don't know yet what they are
	default:
		return nil, fmt.Errorf("Ignoring message with namespace %q", ns)
	}

	// Return success
	return nil, nil
}

// return a new unique message counter
func (this *Channel) nextMsg() int {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.msg = (this.msg + 1) % 100000 // Cycle messages from 1 to 99999
	return this.msg
}

// process heartbeat messages
func (this *Channel) recvHeartbeat(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader

	// Decode the request
	payload := []byte(message.GetPayloadUtf8())
	if err := json.Unmarshal(payload, &header); err != nil {
		return nil, err
	}

	// Emit the payload for debugging
	evt := NewPayloadState(this.key, header.RequestId, payload)
	if err := this.Publisher.Emit(evt, false); err != nil {
		return nil, err
	}

	// Return reply
	switch header.Type {
	case "PING":
		this.ping = time.Now()
		payload := &PayloadHeader{Type: "PONG", RequestId: -1}
		src := message.GetSourceId()
		dst := message.GetDestinationId()
		ns := message.GetNamespace()
		return this.encode(dst, src, ns, payload)
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}
}

// process receiver messages
func (this *Channel) rcvReceiver(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader

	payload := []byte(message.GetPayloadUtf8())
	if err := json.Unmarshal(payload, &header); err != nil {
		return nil, err
	}

	switch header.Type {
	case "RECEIVER_STATUS":
		var status ReceiverStatusResponse

		// Decode the response
		if err := json.Unmarshal(payload, &status); err != nil {
			return nil, fmt.Errorf("RECEIVER_STATUS: %w", err)
		}

		// Emit the volume and applications
		evt := NewAppState(this.key, header.RequestId, payload, status.Status.Volume, status.Status.Apps...)
		if err := this.Publisher.Emit(evt, false); err != nil {
			return nil, err
		}
	case "INVALID_REQUEST", "LAUNCH_ERROR":
		err := new(ErrorResponse)

		// Decode the response
		if err := json.Unmarshal(payload, &err); err != nil {
			return nil, fmt.Errorf("%v: %w", header.Type, err)
		}

		// Emit error
		evt := NewErrorState(this.key, header.RequestId, payload, err)
		if err := this.Publisher.Emit(evt, false); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}

	// Return success
	return nil, nil
}

// process connection messages
func (this *Channel) rcvConnection(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader

	payload := []byte(message.GetPayloadUtf8())
	if err := json.Unmarshal(payload, &header); err != nil {
		return nil, err
	}

	switch header.Type {
	case "CLOSE":
		// Emit the payload for debugging
		evt := NewPayloadState(this.key, header.RequestId, payload)
		if err := this.Publisher.Emit(evt, false); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}

	// Return success
	return nil, nil
}

// process media messages
func (this *Channel) rcvMedia(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader

	payload := []byte(message.GetPayloadUtf8())
	if err := json.Unmarshal(payload, &header); err != nil {
		return nil, err
	}

	switch header.Type {
	case "MEDIA_STATUS":
		var status MediaStatusResponse
		if err := json.Unmarshal(payload, &status); err != nil {
			return nil, fmt.Errorf("MEDIA_STATUS: %w", err)
		}

		// Emit the media state
		evt := NewMediaState(this.key, header.RequestId, payload, status.Status...)
		if err := this.Publisher.Emit(evt, false); err != nil {
			return nil, err
		}
	case "INVALID_REQUEST":
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(header.Type)
	case "LOAD_FAILED", "ERROR":
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(header.Type)
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}

	// Return success
	return nil, nil
}
