// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"golang.org/x/sys/unix"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	InotifyMask uint64
	InotifyEvt  struct {
		unix.InotifyEvent
		path string
	}
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	IN_NONE          InotifyMask = 0
	IN_ACCESS        InotifyMask = (unix.IN_ACCESS)        // File was accessed (read)
	IN_ATTRIB        InotifyMask = (unix.IN_ATTRIB)        // Metadata changed
	IN_CLOSE_WRITE   InotifyMask = (unix.IN_CLOSE_WRITE)   // File opened for writing was closed
	IN_CLOSE_NOWRITE InotifyMask = (unix.IN_CLOSE_NOWRITE) // File not opened for writing was closed
	IN_CREATE        InotifyMask = (unix.IN_CREATE)        // File/directory created in watched directory
	IN_DELETE        InotifyMask = (unix.IN_DELETE)        // File/directory deleted from watched directory
	IN_DELETE_SELF   InotifyMask = (unix.IN_DELETE_SELF)   // Watched file/directory was itself deleted.
	IN_MODIFY        InotifyMask = (unix.IN_MODIFY)        // File was modified
	IN_MOVE_SELF     InotifyMask = (unix.IN_MOVE_SELF)     // Watched file/directory was itself moved.
	IN_MOVED_FROM    InotifyMask = (unix.IN_MOVED_FROM)    // File moved out of watched directory
	IN_MOVED_TO      InotifyMask = (unix.IN_MOVED_TO)      // File moved into watched directory
	IN_OPEN          InotifyMask = (unix.IN_OPEN)          // File was opened
	IN_DONT_FOLLOW   InotifyMask = (unix.IN_DONT_FOLLOW)   // Don't dereference pathname if it is a symbolic link
	IN_EXCL_UNLINK   InotifyMask = (unix.IN_EXCL_UNLINK)   // events are not generated for children after they have been unlinked
	IN_MASK_ADD      InotifyMask = (unix.IN_MASK_ADD)      // Add (OR) events to watch mask for this pathname if it already exists
	IN_ONESHOT       InotifyMask = (unix.IN_ONESHOT)       // Monitor pathname for one event, then remove from watch list.
	IN_ONLYDIR       InotifyMask = (unix.IN_ONLYDIR)       // Only watch pathname if it is a directory
	IN_IGNORED       InotifyMask = (unix.IN_IGNORED)       // Watch was removed explicitly
	IN_ISDIR         InotifyMask = (unix.IN_ISDIR)         // Subject of this event is a directory
	IN_Q_OVERFLOW    InotifyMask = (unix.IN_Q_OVERFLOW)    // Event queue overflowed
	IN_UNMOUNT       InotifyMask = (unix.IN_UNMOUNT)       // File system containing watched object was unmounted
	IN_DEFAULT       InotifyMask = IN_ACCESS | IN_ATTRIB | IN_CREATE | IN_DELETE | IN_DELETE_SELF | IN_MODIFY | IN_MOVE_SELF | IN_MOVED_FROM | IN_MOVED_TO | IN_ISDIR | IN_UNMOUNT
	IN_MAX           InotifyMask = IN_ONESHOT
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	// From https://github.com/rjeczalik/notify
	b      = make([]byte, unix.SizeofInotifyEvent+unix.PathMax+1)
	buffer = bytes.NewBuffer(b)
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func InotifyInit() (*os.File, error) {
	// Create inotify fd
	if fd, err := unix.InotifyInit(); err != nil {
		return nil, err
	} else if fh := os.NewFile(uintptr(fd), ""); fh == nil {
		return nil, gopi.ErrInternalAppError
	} else {
		return fh, nil
	}
}

func InotifyAddWatch(fd uintptr, path string, mask InotifyMask) (uint32, error) {
	if watch, err := unix.InotifyAddWatch(int(fd), path, uint32(mask)); err != nil {
		return 0, err
	} else {
		return uint32(watch), nil
	}
}

func InotifyRemoveWatch(fd uintptr, watch uint32) error {
	if _, err := unix.InotifyRmWatch(int(fd), watch); err != nil {
		return err
	} else {
		return nil
	}
}

func InotifyRead(fd uintptr) (*InotifyEvt, error) {
	if n, err := unix.Read(int(fd), buffer.Bytes()[:]); err != nil {
		return nil, err
	} else if n >= unix.SizeofInotifyEvent {
		var evt InotifyEvt
		if err := binary.Read(buffer, binary.LittleEndian, &evt.InotifyEvent); err != nil {
			return nil, err
		}
		if evt.InotifyEvent.Len > 0 {
			// Read until null
			if path, err := buffer.ReadBytes(0); err == nil {
				evt.path = string(path[0 : len(path)-1])
			}
		}
		return &evt, nil
	} else {
		return nil, gopi.ErrInternalAppError
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this InotifyEvt) Watch() uint32 {
	return uint32(this.InotifyEvent.Wd)
}

func (this InotifyEvt) Mask() InotifyMask {
	return InotifyMask(this.InotifyEvent.Mask)
}

func (this InotifyEvt) Cookie() uint32 {
	return this.InotifyEvent.Cookie
}

func (this InotifyEvt) Len() uint32 {
	return this.InotifyEvent.Len
}

func (this InotifyEvt) Path() string {
	return this.path
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this InotifyEvt) String() string {
	str := " watch=" + fmt.Sprint(this.Watch()) +
		" mask=" + fmt.Sprint(this.Mask())
	if this.Cookie() > 0 {
		str += " cookie=" + fmt.Sprint(this.Cookie())
	}
	if this.Len() > 0 {
		str += " len=" + fmt.Sprint(this.Len())
	}
	str += " path=" + strconv.Quote(this.path)

	return "<InotifyEvt" + str + ">"
}

func (v InotifyMask) String() string {
	if v == IN_NONE {
		return v.StringFlag()
	}
	str := ""
	for f := InotifyMask(1); f <= IN_MAX; f <<= 1 {
		if v&f == f {
			str += f.StringFlag() + "|"
		}
	}
	if str == "" {
		return fmt.Sprintf("%08X", uint32(v))
	} else {
		return strings.TrimSuffix(str, "|")
	}
}

func (v InotifyMask) StringFlag() string {
	switch v {
	case IN_NONE:
		return "IN_NONE"
	case IN_ACCESS:
		return "IN_ACCESS"
	case IN_ATTRIB:
		return "IN_ATTRIB"
	case IN_CLOSE_WRITE:
		return "IN_CLOSE_WRITE"
	case IN_CLOSE_NOWRITE:
		return "IN_CLOSE_NOWRITE"
	case IN_CREATE:
		return "IN_CREATE"
	case IN_DELETE:
		return "IN_DELETE"
	case IN_DELETE_SELF:
		return "IN_DELETE_SELF"
	case IN_MODIFY:
		return "IN_MODIFY"
	case IN_MOVE_SELF:
		return "IN_MOVE_SELF"
	case IN_MOVED_FROM:
		return "IN_MOVED_FROM"
	case IN_MOVED_TO:
		return "IN_MOVED_TO"
	case IN_OPEN:
		return "IN_OPEN"
	case IN_DONT_FOLLOW:
		return "IN_DONT_FOLLOW"
	case IN_EXCL_UNLINK:
		return "IN_EXCL_UNLINK"
	case IN_MASK_ADD:
		return "IN_MASK_ADD"
	case IN_ONESHOT:
		return "IN_ONESHOT"
	case IN_ONLYDIR:
		return "IN_ONLYDIR"
	case IN_IGNORED:
		return "IN_IGNORED"
	case IN_ISDIR:
		return "IN_ISDIR"
	case IN_Q_OVERFLOW:
		return "IN_Q_OVERFLOW"
	case IN_UNMOUNT:
		return "IN_UNMOUNT"
	default:
		return "[?? Invalid InotifyMask value]"
	}
}
