package gopi

import "strings"

/////////////////////////////////////////////////////////////////////
// TYPES

type (
	// Read and Write flags
	FilePollFlags uint

	// FilePollFunc is the handler for file polling
	FilePollFunc func(uintptr, FilePollFlags)
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type FilePoll interface {
	// Watch a file descriptor for changes
	Watch(uintptr, FilePollFlags, FilePollFunc) error

	// Unwatch a file descriptor
	Unwatch(uintptr) error
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FILEPOLL_FLAG_READ  FilePollFlags = (1 << iota) // File descriptor ready for reading
	FILEPOLL_FLAG_WRITE                             // File descriptor ready for writing
	FILEPOLL_FLAG_NONE  FilePollFlags = 0
	FILEPOLL_FLAG_MIN   FilePollFlags = FILEPOLL_FLAG_READ
	FILEPOLL_FLAG_MAX   FilePollFlags = FILEPOLL_FLAG_WRITE
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f FilePollFlags) String() string {
	str := ""
	if f == FILEPOLL_FLAG_NONE {
		return f.StringFlag()
	}
	for v := FILEPOLL_FLAG_MIN; v <= FILEPOLL_FLAG_MAX; v <<= 1 {
		if f&v == v {
			str += v.StringFlag() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (f FilePollFlags) StringFlag() string {
	switch f {
	case FILEPOLL_FLAG_NONE:
		return "FILEPOLL_FLAG_NONE"
	case FILEPOLL_FLAG_READ:
		return "FILEPOLL_FLAG_READ"
	case FILEPOLL_FLAG_WRITE:
		return "FILEPOLL_FLAG_WRITE"
	default:
		return "[?? Invalid FilePollFlags value]"
	}
}
