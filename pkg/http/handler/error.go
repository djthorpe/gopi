package handler

import (
	"net/http"
	"strings"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type err struct {
	req    *http.Request
	code   int
	reason string
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func Error(req *http.Request, code int, reason ...string) error {
	this := new(err)
	this.code = code
	if len(reason) == 0 {
		this.reason = http.StatusText(code)
	} else {
		this.reason = strings.Join(reason, ": ")
	}
	return this
}

/////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *err) Code() int {
	return this.code
}

func (this *err) Error() string {
	if this.req != nil && this.req.URL != nil {
		return this.req.URL.String() + ": " + this.reason
	} else {
		return this.reason
	}
}
