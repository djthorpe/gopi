/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// FilePoll emits READ and WRITE events for files
type FilePoll interface {
	Watch(fd uintptr) error
}
