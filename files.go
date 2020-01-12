/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// READ and WRITE flags
	FilePollFlags uint

	// FilePollFunc is the handler for file polling
	FilePollFunc func(uintptr, FilePollFlags)
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

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FILEPOLL_FLAG_NONE FilePollFlags = (1 << iota)
	FILEPOLL_FLAG_READ
	FILEPOLL_FLAG_WRITE
)
