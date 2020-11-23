package ffmpeg

import (
	"regexp"
	"strconv"
)

var (
	reTrackDisc = regexp.MustCompile("^(\\d+)/(\\d+)$")
)

func ParseTrackDisc(value string) (uint, uint) {
	// parse nn/mm
	if nm := reTrackDisc.FindStringSubmatch(value); len(nm) == 3 {
		if n, err := strconv.ParseUint(nm[1], 0, 64); err != nil {
			return 0, 0
		} else if m, err := strconv.ParseUint(nm[2], 0, 64); err != nil {
			return 0, 0
		} else {
			return uint(n), uint(m)
		}
	} else if n, err := strconv.ParseUint(value, 0, 64); err != nil {
		return 0, 0
	} else {
		return uint(n), 0
	}
}
