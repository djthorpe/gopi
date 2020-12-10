package chromaprint

import (
	"time"
)

type Manager interface {
	// NewStream will create a new stream object for calculating the fingerprint
	NewStream(rate, channels int) (Stream, error)

	// Close releases stream resources
	Close(Stream) error
}

type Stream interface {
	// Methods

	// Write signed data into the stream, assuming little endian order
	Write(data []int16) error

	// GetFingerprint calculates the fingerprint from the streamed data
	GetFingerprint() (string, error)

	// Properties
	// Duration returns the current duration of the streamed data
	Duration() time.Duration

	// Channels returns the number of audio channels the stream represents
	Channels() int

	// Rate returns the samples per second (aka Sample rate)
	Rate() int
}
