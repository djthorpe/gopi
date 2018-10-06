package errors

import (
	"fmt"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// CompoundError can contain one or more errors
type CompoundError struct {
	errs []error
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Add one or more errors onto the array of errors. If any error
// is a CoumpoundError, then add the errors individually
func (this *CompoundError) Add(e ...error) {
	if this.errs == nil {
		this.errs = make([]error, 0, len(e))
	}
	for _, err := range e {
		switch err.(type) {
		case (*CompoundError):
			if len(err.(*CompoundError).errs) > 0 {
				this.errs = append(this.errs, err.(*CompoundError).errs...)
			}
		default:
			this.errs = append(this.errs, err)
		}
	}
}

// Success returns true if no errors appended
func (this *CompoundError) Success() bool {
	return len(this.errs) == 0
}

// One returns the first error if there is only one, or
// else returns nil
func (this *CompoundError) One() error {
	if len(this.errs) == 1 {
		return this.errs[0]
	} else {
		return nil
	}
}

// ErrorOrSelf returns nil, the first error or self
// if there is more than one error
func (this *CompoundError) ErrorOrSelf() error {
	if len(this.errs) == 1 {
		return this.errs[0]
	} else if len(this.errs) == 0 {
		return nil
	} else {
		return this
	}
}

// Error satisfies the error interface
func (this *CompoundError) Error() string {
	if len(this.errs) == 0 {
		return ""
	}
	if len(this.errs) == 1 {
		return this.errs[0].Error()
	}
	errs := ""
	for i, e := range this.errs {
		errs += fmt.Sprintf("Error[%v of %v] %v\n", i+1, len(this.errs), e.Error())
	}
	return strings.Trim(errs, "\n")
}
