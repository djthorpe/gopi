/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package util /* import "github.com/djthorpe/gopi/util" */

import (
	"syscall"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	R_OK = 4
	W_OK = 2
	X_OK = 1
)

////////////////////////////////////////////////////////////////////////////////
// File utility methods

// Returns boolean value which indicates if a file is readable by current
// user
func isReadableFileAtPath(path string) error {
	return syscall.Access(path, R_OK)
}

// Returns boolean value which indicates if a file is writable by current
// user
func isWritableFileAtPath(path string) error {
	return syscall.Access(path, W_OK)
}

// Returns boolean value which indicates if a file is executable by current
// user
func isExecutableFileAtPath(path string) error {
	return syscall.Access(path, X_OK)
}

