package chromecast

// Ref: https://github.com/vishen/go-chromecast/

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Payload interface {
	WithId(id int) Payload
}

type PayloadHeader struct {
	Type      string `json:"type"`
	RequestId int    `json:"requestId,omitempty"`
}

type SetVolumeRequest struct {
	PayloadHeader
	Volume `json:"volume"`
}

type LaunchAppRequest struct {
	PayloadHeader
	AppId string `json:"appId"`
}

type LoadMediaRequest struct {
	PayloadHeader
	Media       MediaItem `json:"media"`
	CurrentTime int       `json:"currentTime,omitempty"`
	Autoplay    bool      `json:"autoplay,omitempty"`
	ResumeState string    `json:"resumeState,omitempty"`
}

type MediaRequest struct {
	PayloadHeader
	MediaSessionId int     `json:"mediaSessionId"`
	CurrentTime    float32 `json:"currentTime,omitempty"`
	RelativeTime   float32 `json:"relativeTime,omitempty"`
	ResumeState    string  `json:"resumeState,omitempty"`
}

type LoadQueueRequest struct {
	PayloadHeader
	RepeatMode string          `json:"repeatMode"`
	Items      []LoadQueueItem `json:"items"`
}

type LoadQueueItem struct {
	PayloadHeader
	Media            MediaItem `json:"media"`
	Autoplay         bool      `json:"autoplay"`
	PlaybackDuration uint      `json:"playbackDuration"`
}

type QueueUpdate struct {
	PayloadHeader
	Media            MediaItem `json:"media"`
	Autoplay         bool      `json:"autoplay"`
	PlaybackDuration uint      `json:"playbackDuration"`
}

type ReceiverStatusResponse struct {
	PayloadHeader
	Status struct {
		Apps   []App `json:"applications"`
		Volume `json:"volume"`
	} `json:"status"`
}

type MediaStatusResponse struct {
	PayloadHeader
	Status []Media `json:"status"`
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *PayloadHeader) WithId(id int) Payload {
	this.RequestId = id
	return this
}

func (this *SetVolumeRequest) WithId(id int) Payload {
	this.PayloadHeader.RequestId = id
	return this
}

func (this *LaunchAppRequest) WithId(id int) Payload {
	this.PayloadHeader.RequestId = id
	return this
}

func (this *LoadMediaRequest) WithId(id int) Payload {
	this.PayloadHeader.RequestId = id
	return this
}

func (this *MediaRequest) WithId(id int) Payload {
	this.PayloadHeader.RequestId = id
	return this
}
