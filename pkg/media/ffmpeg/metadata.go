// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type metadata struct {
	dict *ffmpeg.AVDictionary
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewMetadata(dict *ffmpeg.AVDictionary) *metadata {
	if dict == nil {
		return nil
	} else {
		return &metadata{dict}
	}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *metadata) Keys() []gopi.MediaKey {
	keys := make([]gopi.MediaKey, 0, this.dict.Count())
	entry := this.dict.Get("", nil, ffmpeg.AV_DICT_IGNORE_SUFFIX)
	for entry != nil {
		keys = append(keys, gopi.MediaKey(entry.Key()))
		entry = this.dict.Get("", entry, ffmpeg.AV_DICT_IGNORE_SUFFIX)
	}
	return keys
}

func (this *metadata) Value(key gopi.MediaKey) interface{} {
	if entry := this.dict.Get(string(key), nil, ffmpeg.AV_DICT_IGNORE_SUFFIX); entry == nil {
		return nil
	} else if key == gopi.MEDIA_KEY_COMPILATION {
		if value, err := strconv.ParseInt(entry.Value(), 0, 32); err == nil {
			return value != 0
		} else {
			return nil
		}
	} else if key == gopi.MEDIA_KEY_GAPLESS_PLAYBACK {
		if value, err := strconv.ParseInt(entry.Value(), 0, 32); err == nil {
			return value != 0
		} else {
			return nil
		}
	} else if key == gopi.MEDIA_KEY_TRACK || key == gopi.MEDIA_KEY_DISC {
		n, _ := ParseTrackDisc(entry.Value())
		if n != 0 {
			return n
		} else {
			return nil
		}
	} else if key == gopi.MEDIA_KEY_YEAR {
		if value, err := strconv.ParseUint(entry.Value(), 0, 32); err == nil {
			return uint(value)
		} else {
			return nil
		}
	} else {
		return entry.Value()
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *metadata) String() string {
	str := "<ffmpeg.metadata"
	if keys := this.Keys(); len(keys) > 0 {
		str += " keys="
		for i, key := range keys {
			if i > 0 {
				str += ","
			}
			str += fmt.Sprint(key)
		}
	}
	return str + ">"
}
