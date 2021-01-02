package googlecast

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	pb "github.com/djthorpe/gopi/v3/pkg/rpc/castchannel"
	proto "github.com/golang/protobuf/proto"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type channel struct {
	sync.RWMutex
	msg  int
	key  string
	ch   chan<- state
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

func (this *channel) Init(key string, ch chan<- state) {
	this.key = key
	this.ch = ch
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Generate a connect message
func (this *channel) Connect() (int, []byte, error) {
	payload := &PayloadHeader{Type: "CONNECT"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

// Generate a disconnect message
func (this *channel) Disconnect() (int, []byte, error) {
	payload := &PayloadHeader{Type: "CLOSE"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

// Get status
func (this *channel) GetStatus() (int, []byte, error) {
	payload := &PayloadHeader{Type: "GET_STATUS"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

// Launch application
func (this *channel) LaunchAppWithId(appId string) (int, []byte, error) {
	payload := &LaunchAppRequest{PayloadHeader{Type: "LAUNCH"}, appId}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

// Volume
func (this *channel) SetVolume(v Volume) (int, []byte, error) {
	payload := &SetVolumeRequest{PayloadHeader{Type: "SET_VOLUME"}, v}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

// Mute
func (this *channel) SetMuted(muted bool) (int, []byte, error) {
	payload := &SetVolumeRequest{PayloadHeader{Type: "SET_VOLUME"}, Volume{Muted: muted}}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

// Connect Media
func (this *channel) ConnectMedia(transportId string) (int, []byte, error) {
	payload := &PayloadHeader{Type: "CONNECT"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

// Disconnect Media
func (this *channel) DisconnectMedia(transportId string) (int, []byte, error) {
	payload := &PayloadHeader{Type: "CLOSE"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

// Get Media Status
func (this *channel) GetMediaStatus(transportId string) (int, []byte, error) {
	payload := &PayloadHeader{Type: "GET_STATUS"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_MEDIA, payload.WithId(id))
	return id, data, err
}

// Load
func (this *channel) LoadUrl(transportId string, url, mimetype string, autoplay bool) (int, []byte, error) {
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

// Play
func (this *channel) Play(transportId string, sessionId int, state bool) (int, []byte, error) {
	payload := &MediaRequest{
		MediaSessionId: sessionId,
	}
	switch state {
	case true:
		payload.Type = "PLAY"
	case false:
		payload.Type = "STOP"
	}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_MEDIA, payload.WithId(id))
	return id, data, err
}

// Pause
func (this *channel) Pause(transportId string, sessionId int, state bool) (int, []byte, error) {
	payload := &MediaRequest{
		MediaSessionId: sessionId,
	}
	switch state {
	case false:
		payload.Type = "PLAY"
	case true:
		payload.Type = "PAUSE"
	}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_MEDIA, payload.WithId(id))
	return id, data, err
}

// Seek
func (this *channel) SeekAbs(transportId string, sessionId int, value float32) (int, []byte, error) {
	payload := &MediaRequest{
		PayloadHeader: PayloadHeader{
			Type: "SEEK",
		},
		MediaSessionId: sessionId,
		CurrentTime:    value,
		ResumeState:    "PLAYBACK_START",
	}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_MEDIA, payload.WithId(id))
	return id, data, err
}

// SeekRel
func (this *channel) SeekRel(transportId string, sessionId int, value float32) (int, []byte, error) {
	payload := &MediaRequest{
		PayloadHeader: PayloadHeader{
			Type: "SEEK",
		},
		MediaSessionId: sessionId,
		RelativeTime:   value,
		ResumeState:    "PLAYBACK_START",
	}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, transportId, CAST_NS_MEDIA, payload.WithId(id))
	return id, data, err
}

// Stop
func (this *channel) Stop() (int, []byte, error) {
	payload := &PayloadHeader{
		Type: "STOP",
	}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_RECV, payload.WithId(id))
	return id, data, err
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// encode message and return it
func (this *channel) encode(source, dest, ns string, payload Payload) ([]byte, error) {
	//if debug, err := json.MarshalIndent(payload, "", "  "); err == nil {
	//	fmt.Println(string(debug))
	//}
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
func (this *channel) decode(data []byte) ([]byte, error) {
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
func (this *channel) nextMsg() int {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.msg = (this.msg + 1) % 100000 // Cycle messages from 1 to 99999
	return this.msg
}

// process heartbeat messages
func (this *channel) recvHeartbeat(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
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
func (this *channel) rcvReceiver(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
	switch header.Type {
	case "RECEIVER_STATUS":
		var status ReceiverStatusResponse

		// Decode response
		if err := json.Unmarshal([]byte(message.GetPayloadUtf8()), &status); err != nil {
			return nil, fmt.Errorf("RECEIVER_STATUS: %w", err)
		}

		// Emit the volume and first application (doesn't support more than one)
		if len(status.Status.Apps) == 0 {
			this.ch <- NewState(this.key, header.RequestId, status.Status.Volume, App{})
		} else {
			this.ch <- NewState(this.key, header.RequestId, status.Status.Volume, status.Status.Apps[0])
		}
	case "INVALID_REQUEST", "LAUNCH_ERROR":
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(message.GetPayloadUtf8())
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}

	// Return success
	return nil, nil
}

// process connection messages
func (this *channel) rcvConnection(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
	switch header.Type {
	//case "CLOSE":
	//	this.ch <- NewState(this.key, header.RequestId, "CLOSE")
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}

	// Return success
	//return nil, nil
}

// process media messages
func (this *channel) rcvMedia(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
	switch header.Type {
	case "MEDIA_STATUS":
		var status MediaStatusResponse
		if err := json.Unmarshal([]byte(message.GetPayloadUtf8()), &status); err != nil {
			return nil, fmt.Errorf("MEDIA_STATUS: %w", err)
		}
		// Emit the media items
		if len(status.Status) == 0 {
			this.ch <- NewState(this.key, header.RequestId, Media{})
		} else {
			result := make([]interface{}, len(status.Status))
			for i, media := range status.Status {
				result[i] = media
			}
			this.ch <- NewState(this.key, header.RequestId, result...)
		}
	case "INVALID_REQUEST":
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(message.GetPayloadUtf8())
	case "LOAD_FAILED", "ERROR":
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(message.GetPayloadUtf8())
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}

	// Return success
	return nil, nil

}
