package googlecast

import (
	"fmt"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Media struct {
	MediaSessionId int `json:"mediaSessionId"`

	// PLAYING, BUFFERING, PAUSED, IDLE or UNKNOWN
	PlayerState string `json:"playerState"`

	// CurrentTime in seconds
	CurrentTime float32 `json:"currentTime"`

	IdleReason    string    `json:"idleReason"`
	Volume        Volume    `json:"volume"`
	CurrentItemId int       `json:"currentItemId"`
	LoadingItemId int       `json:"loadingItemId"`
	Media         MediaItem `json:"media"`
}

type MediaItem struct {
	ContentId string `json:"contentId"`

	// Mimetype
	ContentType string `json:"contentType,omitempty"`

	// BUFFERED, LIVE or UNKNOWN
	StreamType string `json:"streamType,omitempty"`

	// Duration in seconds
	Duration float32        `json:"duration,omitempty"`
	Metadata *MediaMetadata `json:"metadata,omitempty"`
}

type MediaMetadata struct {
	MetadataType int          `json:"metadataType,omitempty"`
	Artist       string       `json:"artist,omitempty"`
	Title        string       `json:"title,omitempty"`
	Subtitle     string       `json:"subtitle,omitempty"`
	Images       []MediaImage `json:"images,omitempty"`
	ReleaseDate  string       `json:"releaseDate,omitempty"`
}

type MediaImage struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (m Media) Equals(other Media) bool {
	if m.MediaSessionId != other.MediaSessionId {
		return false
	}
	if m.PlayerState != other.PlayerState {
		return false
	}
	if m.CurrentTime != other.CurrentTime {
		return false
	}
	if m.IdleReason != other.IdleReason {
		return false
	}
	if m.CurrentItemId != other.CurrentItemId {
		return false
	}
	if m.LoadingItemId != other.LoadingItemId {
		return false
	}
	return m.Media.Equals(other.Media)
}

func (m MediaItem) Equals(other MediaItem) bool {
	if m.ContentId != other.ContentId {
		return false
	}
	if m.ContentType != other.ContentType {
		return false
	}
	if m.StreamType != other.StreamType {
		return false
	}
	if m.Duration != other.Duration {
		return false
	}
	return m.Metadata.Equals(other.Metadata)
}

func (m *MediaMetadata) Equals(other *MediaMetadata) bool {
	if other == nil {
		return m == nil
	}
	if m.MetadataType != other.MetadataType {
		return false
	}
	if m.Artist != other.Artist {
		return false
	}
	if m.Title != other.Title {
		return false
	}
	if m.Subtitle != other.Subtitle {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (m Media) String() string {
	var parts string
	if m.PlayerState != "" {
		parts += fmt.Sprintf(" state=%v", strconv.Quote(m.PlayerState))
	}
	if m.IdleReason != "" {
		parts += fmt.Sprintf(" idle_reason=%v", strconv.Quote(m.IdleReason))
	}
	if m.CurrentItemId != 0 {
		parts += fmt.Sprintf(" current_id=%v", m.CurrentItemId)
	}
	if m.CurrentTime != 0 {
		parts += fmt.Sprintf(" current_time=%v", m.CurrentTime)
	}
	if m.LoadingItemId != 0 {
		parts += fmt.Sprintf(" loading_id=%v", m.LoadingItemId)
	}
	if m.Media.ContentId != "" {
		parts += fmt.Sprintf(" %v", m.Media)
	}
	return fmt.Sprintf("<media media_session_id=%v%v>", m.MediaSessionId, parts)
}

func (m MediaItem) String() string {
	var parts string
	if m.ContentType != "" {
		parts += fmt.Sprintf(" content_type=%v", strconv.Quote(m.ContentType))
	}
	if m.StreamType != "" {
		parts += fmt.Sprintf(" stream_type=%v", strconv.Quote(m.StreamType))
	}
	if m.Duration != 0 {
		parts += fmt.Sprintf(" duration=%v", m.Duration)
	}
	if m.Metadata != nil && m.Metadata.MetadataType != 0 {
		parts += fmt.Sprintf(" %v", m.Metadata)
	}
	return fmt.Sprintf("<item content_id=%v%v>", m.ContentId, parts)
}

func (m MediaMetadata) String() string {
	var parts string
	if m.Artist != "" {
		parts += fmt.Sprintf(" artist=%v", strconv.Quote(m.Artist))
	}
	if m.Title != "" {
		parts += fmt.Sprintf(" title=%v", strconv.Quote(m.Title))
	}
	if m.Subtitle != "" {
		parts += fmt.Sprintf(" subtitle=%v", strconv.Quote(m.Subtitle))
	}
	if m.ReleaseDate != "" {
		parts += fmt.Sprintf(" release_date=%v", strconv.Quote(m.ReleaseDate))
	}
	if len(m.Images) > 0 {
		parts += fmt.Sprintf(" images=%v", m.Images)
	}
	return fmt.Sprintf("<metadata type=%v%v>", m.MetadataType, parts)
}

func (m MediaImage) String() string {
	var parts string
	if m.Width != 0 {
		parts += fmt.Sprintf(" w=%v", m.Width)
	}
	if m.Height != 0 {
		parts += fmt.Sprintf(" h=%v", m.Height)
	}
	return fmt.Sprintf("<image url=%v%v>", m.URL, parts)
}
