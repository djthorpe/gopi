/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import "strings"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// READ and WRITE flags
	FilePollFlags uint

	// File system event flags
	FSEventFlags uint

	// FilePollFunc is the handler for file polling
	FilePollFunc func(uintptr, FilePollFlags)

	// FSEventFunc is the handler for filesystem events
	FSEventFunc func(FSEventFlags)
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// FilePoll emits READ and WRITE events for files
type FilePoll interface {
	// Watch for events and call handler with file descriptor
	Watch(uintptr, FilePollFlags, FilePollFunc) error

	// Unwatch stops watching for file events
	Unwatch(uintptr) error

	// Implements gopi.Unit
	Unit
}

// FSEvents emits events when watched files and folders are changed
type FSEvents interface {
	// Watch for events and call handler when events are generated
	Watch(string, FSEventFlags, FSEventFunc) (uint32, error)

	// Unwatch stops watching for file events
	Unwatch(uint32) error

	// Implements gopi.Unit
	Unit
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

const (
	FSEVENT_FLAG_ATTRIB  FSEventFlags = (1 << iota) // Attributes changed
	FSEVENT_FLAG_CREATE                             // File created
	FSEVENT_FLAG_DELETE                             // File deleted
	FSEVENT_FLAG_MODIFY                             // File modified
	FSEVENT_FLAG_MOVE                               // File moved
	FSEVENT_FLAG_UNMOUNT                            // Watched filesystem was unmounted
	FSEVENT_FLAG_ALL                  = FSEVENT_FLAG_ATTRIB | FSEVENT_FLAG_CREATE | FSEVENT_FLAG_DELETE | FSEVENT_FLAG_MODIFY | FSEVENT_FLAG_MOVE | FSEVENT_FLAG_UNMOUNT
	FSEVENT_FLAG_NONE    FSEventFlags = 0
	FSEVENT_FLAG_MIN     FSEventFlags = FSEVENT_FLAG_ATTRIB
	FSEVENT_FLAG_MAX     FSEventFlags = FSEVENT_FLAG_UNMOUNT
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

func (f FSEventFlags) String() string {
	str := ""
	if f == FSEVENT_FLAG_NONE {
		return f.StringFlag()
	}
	for v := FSEVENT_FLAG_MIN; v <= FSEVENT_FLAG_MAX; v <<= 1 {
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

func (f FSEventFlags) StringFlag() string {
	switch f {
	case FSEVENT_FLAG_NONE:
		return "FSEVENT_FLAG_NONE"
	case FSEVENT_FLAG_ATTRIB:
		return "FSEVENT_FLAG_ATTRIB"
	case FSEVENT_FLAG_CREATE:
		return "FSEVENT_FLAG_CREATE"
	case FSEVENT_FLAG_DELETE:
		return "FSEVENT_FLAG_DELETE"
	case FSEVENT_FLAG_MODIFY:
		return "FSEVENT_FLAG_MODIFY"
	case FSEVENT_FLAG_MOVE:
		return "FSEVENT_FLAG_MOVE"
	case FSEVENT_FLAG_UNMOUNT:
		return "FSEVENT_FLAG_UNMOUNT"
	default:
		return "[?? Invalid FSEventFlags value]"
	}
}
