package mdns

import "fmt"

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	qtypeMap = map[uint16]string{
		1:  "A",
		2:  "NS",
		5:  "CNAME",
		6:  "SOA",
		12: "PTR",
		15: "MX",
		16: "TXT",
		28: "AAAA",
		33: "SRV",
		99: "SPF",
	}
)

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func qTypeString(q uint16) string {
	if str, exists := qtypeMap[q]; exists {
		return str
	} else {
		return fmt.Sprint(q)
	}
}
