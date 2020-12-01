package main

import (
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

func FormatCodes(codes []gopi.KeyCode) string {
	str := ""
	for i, code := range codes {
		if i > 0 {
			str += ", "
		}
		str += fmt.Sprint(code)
	}
	return str
}
