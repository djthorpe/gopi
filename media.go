package gopi

import (
	"net/url"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type MediaKey string
type MediaFlag uint
type DecodeIteratorFunc func(MediaDecodeContext, MediaPacket) error
type DecodeFrameIteratorFunc func(MediaFrame) error

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// MediaManager for media file management
type MediaManager interface {
	OpenFile(path string) (Media, error) // Open a media file
	Close(Media) error                   // Close a media object
}

// Media is an input or output
type Media interface {
	// Properties
	URL() *url.URL           // Return URL for the media location
	Metadata() MediaMetadata // Return metadata
	Flags() MediaFlag        // Return flags
	Streams() []MediaStream  // Return streams

	// DecodeIterator loops over selected streams from media object
	DecodeIterator([]int, DecodeIteratorFunc) error

	// DecodeFrameIterator loops over decoded frames from media stream
	DecodeFrameIterator(MediaDecodeContext, MediaPacket, DecodeFrameIteratorFunc) error
}

// MediaMetadata are key value pairs for a media object
type MediaMetadata interface {
	Keys() []MediaKey           // Return all existing keys
	Value(MediaKey) interface{} // Return value for key, or nil
}

// MediaStream is a stream of packets from a media object
type MediaStream interface {
	Index() int        // Stream index
	Flags() MediaFlag  // Flags for the stream (Audio, Video, etc)
	Codec() MediaCodec // Return codec and parameters
}

// MediaCodec is the codec and parameters
type MediaCodec interface {
	Flags() MediaFlag // Flags for the codec (Audio, Video, etc)
}

// MediaPacket is a packet of data from a stream
type MediaPacket interface {
	Size() int
	Bytes() []byte
	Stream() int
}

// MediaFrame is a decoded audio or video frame
type MediaFrame interface {
}

// MediaDecodeContext provides packet data and streams for decoding
// frames of data
type MediaDecodeContext interface {
	Stream() MediaStream // Origin of the packet
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MEDIA_FLAG_ALBUM             MediaFlag = (1 << iota) // Is part of an album
	MEDIA_FLAG_ALBUM_TRACK                               // Is an album track
	MEDIA_FLAG_ALBUM_COMPILATION                         // Album is a compliation
	MEDIA_FLAG_TVSHOW                                    // Is part of a TV Show
	MEDIA_FLAG_TVSHOW_EPISODE                            // Is a TV Show episode
	MEDIA_FLAG_FILE                                      // Is a file
	MEDIA_FLAG_VIDEO                                     // Contains video
	MEDIA_FLAG_AUDIO                                     // Contains audio
	MEDIA_FLAG_DATA                                      // Contains data stream
	MEDIA_FLAG_ATTACHMENT                                // Contains attachment
	MEDIA_FLAG_ARTWORK                                   // Contains artwork
	MEDIA_FLAG_CAPTIONS                                  // Contains captions
	MEDIA_FLAG_NONE              MediaFlag = 0
	MEDIA_FLAG_MIN                         = MEDIA_FLAG_ALBUM
	MEDIA_FLAG_MAX                         = MEDIA_FLAG_CAPTIONS
)

const (
	MEDIA_KEY_BRAND_MAJOR      MediaKey = "major_brand"       // string
	MEDIA_KEY_BRAND_COMPATIBLE MediaKey = "compatible_brands" // string
	MEDIA_KEY_CREATED          MediaKey = "creation_time"     // time.Time
	MEDIA_KEY_ENCODER          MediaKey = "encoder"           // string
	MEDIA_KEY_ALBUM            MediaKey = "album"             // string
	MEDIA_KEY_ALBUM_ARTIST     MediaKey = "artist"            // string
	MEDIA_KEY_COMMENT          MediaKey = "comment"           // string
	MEDIA_KEY_COMPOSER         MediaKey = "composer"          // string
	MEDIA_KEY_COPYRIGHT        MediaKey = "copyright"         // string
	MEDIA_KEY_YEAR             MediaKey = "date"              // uint
	MEDIA_KEY_DISC             MediaKey = "disc"              // uint
	MEDIA_KEY_ENCODED_BY       MediaKey = "encoded_by"        // string
	MEDIA_KEY_FILENAME         MediaKey = "filename"          // string
	MEDIA_KEY_GENRE            MediaKey = "genre"             // string
	MEDIA_KEY_LANGUAGE         MediaKey = "language"          // string
	MEDIA_KEY_PERFORMER        MediaKey = "performer"         // string
	MEDIA_KEY_PUBLISHER        MediaKey = "publisher"         // string
	MEDIA_KEY_SERVICE_NAME     MediaKey = "service_name"      // string
	MEDIA_KEY_SERVICE_PROVIDER MediaKey = "service_provider"  // string
	MEDIA_KEY_TITLE            MediaKey = "title"             // string
	MEDIA_KEY_TRACK            MediaKey = "track"             // uint
	MEDIA_KEY_VERSION_MAJOR    MediaKey = "major_version"     // string
	MEDIA_KEY_VERSION_MINOR    MediaKey = "minor_version"     // string
	MEDIA_KEY_SHOW             MediaKey = "show"              // string
	MEDIA_KEY_SEASON           MediaKey = "season_number"     // uint
	MEDIA_KEY_EPISODE_SORT     MediaKey = "episode_sort"      // string
	MEDIA_KEY_EPISODE_ID       MediaKey = "episode_id"        // uint
	MEDIA_KEY_COMPILATION      MediaKey = "compilation"       // bool
	MEDIA_KEY_GAPLESS_PLAYBACK MediaKey = "gapless_playback"  // bool
	MEDIA_KEY_ACCOUNT_ID       MediaKey = "account_id"        // string
	MEDIA_KEY_DESCRIPTION      MediaKey = "description"       // string
	MEDIA_KEY_MEDIA_TYPE       MediaKey = "media_type"        // string
	MEDIA_KEY_PURCHASED        MediaKey = "purchase_date"     // time.Time
	MEDIA_KEY_ALBUM_SORT       MediaKey = "sort_album"        // string
	MEDIA_KEY_ARTIST_SORT      MediaKey = "sort_artist"       // string
	MEDIA_KEY_TITLE_SORT       MediaKey = "sort_name"         // string
	MEDIA_KEY_SYNOPSIS         MediaKey = "synopsis"          // string
	MEDIA_KEY_GROUPING         MediaKey = "grouping"          // string
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f MediaFlag) String() string {
	if f == MEDIA_FLAG_NONE {
		return f.FlagString()
	}
	str := ""
	for v := MEDIA_FLAG_MIN; v <= MEDIA_FLAG_MAX; v <<= 1 {
		if f&v == v {
			str += v.FlagString() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (f MediaFlag) FlagString() string {
	switch f {
	case MEDIA_FLAG_NONE:
		return "MEDIA_FLAG_NONE"
	case MEDIA_FLAG_ALBUM:
		return "MEDIA_FLAG_ALBUM"
	case MEDIA_FLAG_ALBUM_TRACK:
		return "MEDIA_FLAG_ALBUM_TRACK"
	case MEDIA_FLAG_ALBUM_COMPILATION:
		return "MEDIA_FLAG_ALBUM_COMPILATION"
	case MEDIA_FLAG_TVSHOW:
		return "MEDIA_FLAG_TVSHOW"
	case MEDIA_FLAG_TVSHOW_EPISODE:
		return "MEDIA_FLAG_TVSHOW_EPISODE"
	case MEDIA_FLAG_FILE:
		return "MEDIA_FLAG_FILE"
	case MEDIA_FLAG_VIDEO:
		return "MEDIA_FLAG_VIDEO"
	case MEDIA_FLAG_AUDIO:
		return "MEDIA_FLAG_AUDIO"
	case MEDIA_FLAG_DATA:
		return "MEDIA_FLAG_DATA"
	case MEDIA_FLAG_ATTACHMENT:
		return "MEDIA_FLAG_ATTACHMENT"
	case MEDIA_FLAG_ARTWORK:
		return "MEDIA_FLAG_ARTWORK"
	case MEDIA_FLAG_CAPTIONS:
		return "MEDIA_FLAG_CAPTIONS"
	default:
		return "[?? Invalid MediaFlag]"
	}
}
