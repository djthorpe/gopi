/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

import (
	"syscall"
)

const (
	R_OK = 4
	W_OK = 2
	X_OK = 1
)

func isReadablePath(path string) error {
	return syscall.Access(path, R_OK)
}

func isWritablePath(path string) error {
	return syscall.Access(path, W_OK)
}

func isExecutableFileAtPath(path string) error {
	return syscall.Access(path, X_OK)
}
