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

// Add appends an error onto the array of errors
func (this *CompoundError) Add(e error) {
	if this.errs == nil {
		this.errs = make([]error, 0, 1)
	}
	if e != nil {
		this.errs = append(this.errs, e)
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
