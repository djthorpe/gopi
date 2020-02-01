/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package files

import (
	"fmt"
	"strconv"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type fsevent struct {
	source     gopi.FSEvents
	watch      uint32
	root, path string
	flags      gopi.FSEventFlags
	isfolder   bool
	inode      uint64
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewEvent(source gopi.FSEvents, watch uint32, root, path string, flags gopi.FSEventFlags, inode uint64, isfolder bool) gopi.FSEvent {
	return &fsevent{source, watch, root, path, flags, isfolder, inode}
}

////////////////////////////////////////////////////////////////////////////////
// FSEvent implementation

func (*fsevent) Name() string {
	return "gopi.FSEvent"
}

func (this *fsevent) Source() gopi.Unit {
	return this.source
}

func (this *fsevent) Value() interface{} {
	return this.watch
}

func (*fsevent) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

// Root path being watched
func (this *fsevent) Root() string {
	return this.root
}

// Path the event is concerned with, under the root path
func (this *fsevent) Path() string {
	return this.path
}

// Whether the event concerns a folder
func (this *fsevent) IsFolder() bool {
	return this.isfolder
}

// The event flags
func (this *fsevent) Flags() gopi.FSEventFlags {
	return this.flags
}

// The unique inode number for the event path, or zero
func (this *fsevent) Inode() uint64 {
	return this.inode
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *fsevent) String() string {
	str := "<gopi.FSEvent" + " root=" + strconv.Quote(this.root)
	if this.path != "" {
		str += " path=" + strconv.Quote(this.path)
	}
	if this.flags != 0 {
		str += " flags=" + fmt.Sprint(this.flags)
	}
	if this.isfolder {
		str += " isfolder=true"
	}
	if this.inode != 0 {
		str += " inode=" + fmt.Sprint(this.inode)
	}
	return str + ">"
}
