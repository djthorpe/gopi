package googlecast

import (
	"encoding/json"
	"fmt"
	"sync"

	pb "github.com/djthorpe/gopi/v3/pkg/rpc/castchannel"
	proto "github.com/golang/protobuf/proto"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type channel struct {
	sync.RWMutex

	msg int
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
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *channel) Connect() (int, []byte, error) {
	payload := &PayloadHeader{Type: "CONNECT"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

func (this *channel) Disconnect() (int, []byte, error) {
	payload := &PayloadHeader{Type: "CLOSE"}
	id := this.nextMsg()
	data, err := this.encode(CAST_DEFAULT_SENDER, CAST_DEFAULT_RECEIVER, CAST_NS_CONN, payload.WithId(id))
	return id, data, err
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *channel) encode(source, dest, ns string, payload Payload) ([]byte, error) {
	fmt.Println(payload)
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
	default:
		return nil, fmt.Errorf("Ignoring message with namespace %q", ns)
	}
}

func (this *channel) nextMsg() int {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.msg = (this.msg + 1) % 100000 // Cycle messages from 1 to 99999
	return this.msg
}

func (this *channel) recvHeartbeat(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
	switch header.Type {
	case "PING":
		payload := &PayloadHeader{Type: "PONG", RequestId: -1}
		src := message.GetSourceId()
		dst := message.GetDestinationId()
		ns := message.GetNamespace()
		return this.encode(dst, src, ns, payload)
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}
}

func (this *channel) rcvReceiver(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
	switch header.Type {
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}
}

func (this *channel) rcvConnection(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
	switch header.Type {
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}
}

func (this *channel) rcvMedia(message *pb.CastMessage) ([]byte, error) {
	var header PayloadHeader
	if err := json.Unmarshal([]byte(*message.PayloadUtf8), &header); err != nil {
		return nil, err
	}
	switch header.Type {
	default:
		return nil, fmt.Errorf("Ignoring message %q in namespace %q", header.Type, message.GetNamespace())
	}
}
