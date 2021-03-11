/*
	Rotel RS232 Control
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package rotel

import (
	"fmt"
	"strings"
	"sync"

	rotel "github.com/djthorpe/rotel"
	// Frameworks
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	mapSource = make(map[string]rotel.Source)
	mapLock   sync.Mutex
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

func stringToSource(value string) rotel.Source {
	if len(mapSource) == 0 {
		mapLock.Lock()
		defer mapLock.Unlock()
		for source := rotel.Source(0); source <= rotel.ROTEL_SOURCE_MAX; source++ {
			str := sourceToString(source)
			mapSource[str] = source
		}
	}
	if src, exists := mapSource[value]; exists {
		return src
	} else {
		return rotel.ROTEL_SOURCE_NONE
	}
}

func sourceToString(value rotel.Source) string {
	str := fmt.Sprint(value)
	if strings.HasPrefix(str, "ROTEL_SOURCE_") == false {
		return ""
	} else if str = strings.TrimPrefix(str, "ROTEL_SOURCE_"); str == "NONE" || str == "OTHER" {
		return ""
	} else {
		return strings.ToLower(str)
	}
}
