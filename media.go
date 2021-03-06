package gopi

import (
	"context"
	"image"
	"io"
	"net/url"
	"strings"
)

/*
	This file contains definitions for media devices:

	* Video and Audio encoding and decoding
	* Input and output media devices
	* DVB tuning and decoding (experimental)

	There are aditional interfaces for audio and graphics elsewhere
*/

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MediaKey                string
	MediaFlag               uint64
	DecodeIteratorFunc      func(MediaDecodeContext, MediaPacket) error
	DecodeFrameIteratorFunc func(MediaFrame) error
)

////////////////////////////////////////////////////////////////////////////////
// MEDIA MANAGER

// MediaManager for media file management
type MediaManager interface {
	// OpenFile opens a local media file
	OpenFile(path string) (MediaInput, error)

	// OpenURL opens a network-based stream
	OpenURL(url *url.URL) (MediaInput, error)

	// CreateFile creates a local media file for output
	CreateFile(path string) (MediaOutput, error)

	// Close will release resources and close a media object
	Close(Media) error

	// ListCodecs enumerates codecs for a specific name and/or
	// audio, video, encode and decode. By default (empty name and
	// MediaFlag) lists all codecs
	ListCodecs(string, MediaFlag) []MediaCodec

	// Create an audio profile with format, sample rate, channels and layout
	AudioProfile(AudioFormat, uint, AudioChannelLayout) MediaProfile
}

////////////////////////////////////////////////////////////////////////////////
// MEDIA OBJECTS

// Media is an input or output
type Media interface {
	URL() *url.URL                  // Return URL for the media location
	Metadata() MediaMetadata        // Return metadata
	Flags() MediaFlag               // Return flags
	Streams() []MediaStream         // Return streams
	StreamForIndex(int) MediaStream // Return stream by index
}

// MediaInput represents a source of media
type MediaInput interface {
	Media

	// StreamsForFlag returns array of stream indexes for
	// the best streams to use according to the flags
	StreamsForFlag(MediaFlag) []int

	// Read loops over selected streams from media object, and
	// packets are provided to a Decode function
	Read(context.Context, []int, DecodeIteratorFunc) error

	// DecodeFrameIterator loops over data packets from media stream
	DecodeFrameIterator(MediaDecodeContext, MediaPacket, DecodeFrameIteratorFunc) error
}

// MediaOutput represents a sink for media
type MediaOutput interface {
	Media

	// Write packets to output
	Write(MediaDecodeContext, MediaPacket) error
}

////////////////////////////////////////////////////////////////////////////////
// MEDIA PROFILES

type MediaProfile interface {
	Flags() MediaFlag // Return audio or video profile
}

type MediaAudioProfile interface {
	MediaProfile

	Format() AudioFormat
	SampleRate() uint
	Layout() AudioChannelLayout
}

type MediaVideoProfile interface {
	MediaProfile
}

////////////////////////////////////////////////////////////////////////////////
// MEDIA METADATA AND CODECS

// MediaMetadata are key value pairs for a media object
type MediaMetadata interface {
	Keys() []MediaKey           // Return all existing keys
	Value(MediaKey) interface{} // Return value for key, or nil
}

// MediaCodec is the codec and parameters
type MediaCodec interface {
	// Name returns the unique name for the codec
	Name() string

	// Description returns the long description for the codec
	Description() string

	// Flags for the codec (Audio, Video, Encoder, Decoder)
	Flags() MediaFlag
}

////////////////////////////////////////////////////////////////////////////////
// MEDIA STREAMS, PACKETS AND FRAMES

// MediaStream is a stream of packets from a media object, with a defined
// codec and audio or video profile
type MediaStream interface {
	Index() int            // Stream index
	Flags() MediaFlag      // Flags for the stream (Audio, Video, etc)
	Codec() MediaCodec     // Return codec and parameters
	Profile() MediaProfile // Return audio or video profile for stream
}

// MediaPacket is a packet of data from a stream
type MediaPacket interface {
	Size() int
	Bytes() []byte
	Stream() int
}

// MediaFrame is a decoded audio or video frame
type MediaFrame interface {
	// Implements image interface which can be saved to save frame as bitmap
	image.Image

	// Resample a frame to a specific profile
	Resample(MediaProfile) (MediaFrame, error)

	// Flags for the frame (Audio, Video)
	Flags() MediaFlag
}

////////////////////////////////////////////////////////////////////////////////
// MEDIA ENCODING AND DECODING

// MediaDecodeContext provides packet data and streams for decoding
// frames of data
type MediaDecodeContext interface {
	Stream() MediaStream // Origin of the packet
	Frame() int          // Frame counter
}

////////////////////////////////////////////////////////////////////////////////
// DVB INTERFACES - EXPERIMENTAL

// DVBManager encapsulates methods for DVB reception
type DVBManager interface {
	// Tuners returns all tuners
	Tuners() []DVBTuner

	// ParseTunerParams returns a list of tuner parameters
	ParseTunerParams(r io.Reader) ([]DVBTunerParams, error)

	// Tune to parameters with context, which can be used to cancel
	// or timeout the tuning process
	Tune(context.Context, DVBTuner, DVBTunerParams, DVBTuneCallack) error

	// Close tuner
	Close(DVBTuner) error
}

// DVBTunerCallack is called when tuning has completed, so that
// set-up can begin
type DVBTuneCallack func(DVBContext)

// DVBTuneContext represents information obtained during the tuning
// process
type DVBContext interface {
	//	Services() []DVBService
}

// DVBTunerParams represents tune parameters
type DVBTunerParams interface {
	// Name returns name for tune parameters
	Name() string
}

// DVBTuner represents a tuning device, some hardware
// comtains more than one tuner. Each is represented by
// a unique identifier
type DVBTuner interface {
	// Return tuner identifier
	Id() uint

	// Return tuner name
	Name() string
}

// DVBService represents a service that can be received
type DVBService interface {
	Id() uint16
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
	MEDIA_FLAG_SUBTITLE                                  // Contains subtitles
	MEDIA_FLAG_DATA                                      // Contains data stream
	MEDIA_FLAG_ATTACHMENT                                // Contains attachment
	MEDIA_FLAG_ARTWORK                                   // Contains artwork
	MEDIA_FLAG_CAPTIONS                                  // Contains captions
	MEDIA_FLAG_ENCODER                                   // Is an encoder
	MEDIA_FLAG_DECODER                                   // Is an decoder
	MEDIA_FLAG_NONE              MediaFlag = 0
	MEDIA_FLAG_MIN                         = MEDIA_FLAG_ALBUM
	MEDIA_FLAG_MAX                         = MEDIA_FLAG_DECODER
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
	case MEDIA_FLAG_SUBTITLE:
		return "MEDIA_FLAG_SUBTITLE"
	case MEDIA_FLAG_DATA:
		return "MEDIA_FLAG_DATA"
	case MEDIA_FLAG_ATTACHMENT:
		return "MEDIA_FLAG_ATTACHMENT"
	case MEDIA_FLAG_ARTWORK:
		return "MEDIA_FLAG_ARTWORK"
	case MEDIA_FLAG_CAPTIONS:
		return "MEDIA_FLAG_CAPTIONS"
	case MEDIA_FLAG_ENCODER:
		return "MEDIA_FLAG_ENCODER"
	case MEDIA_FLAG_DECODER:
		return "MEDIA_FLAG_DECODER"
	default:
		return "[?? Invalid MediaFlag]"
	}
}
