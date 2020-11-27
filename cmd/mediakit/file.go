package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/djthorpe/gopi/v3"
)

type media struct {
	URL   *url.URL
	Flags gopi.MediaFlag
	//	Metadata map[gopi.MediaKey]interface{}
	Streams []*stream
}

type stream struct {
	Flags gopi.MediaFlag
	Index int
}

func NewMedia(m gopi.Media) *media {
	this := new(media)
	this.URL = m.URL()
	this.Flags = m.Flags()
	return this
}

func NewStream(s gopi.MediaStream) *stream {
	this := new(stream)
	this.Index = s.Index()
	this.Flags = s.Flags()
	return this
}

func FormatStreams(m *media) [][]string {
	rows := [][]string{}
	for i, s := range m.Streams {
		name := ""
		if i == 0 {
			name = filepath.Base(m.URL.Path)
		}
		rows = append(rows, FormatStream(name, s))
	}
	return rows
}

func FormatStream(name string, stream *stream) []string {
	return []string{
		name,
		fmt.Sprint(stream.Index),
		FormatFlags(stream.Flags),
	}
}

func FormatFlags(flags gopi.MediaFlag) string {
	str := ""
	for _, flag := range strings.Split(fmt.Sprint(flags), "|") {
		flag := strings.TrimPrefix(flag, "MEDIA_FLAG_")
		str += flag
	}
	return str
}

/*

func FormatArtists(artists []string) string {
	str := ""
	for i, artist := range artists {
		if i > 0 {
			str += ", "
		}
		str += strconv.Quote(artist)
	}
	return str
}

func FormatMetadata(metadata map[gopi.MediaKey]interface{}) string {
	str := ""
	for k, v := range metadata {
		switch v.(type) {
		case string:
			str += fmt.Sprintf(" %s=%q", k, v.(string))
		default:
			str += fmt.Sprintf(" %s=%v", k, v)
		}
	}
	return strings.TrimSpace(str)
}

func (this *album) Artists() []string {
	keys := make(map[string]bool, len(this.files))
	for _, file := range this.files {
		key, ok := file.Metadata[gopi.MEDIA_KEY_ALBUM_ARTIST].(string)
		if ok == false || len(key) < 2 {
			continue
		}
		keys[key] = true
	}
	artists := []string{}
	for k := range keys {
		artists = append(artists, k)
	}
	return artists
}

func (this *file) Filename() string {
	disc, _ := this.Metadata[gopi.MEDIA_KEY_DISC].(uint)
	track, _ := this.Metadata[gopi.MEDIA_KEY_TRACK].(uint)
	title, _ := this.Metadata[gopi.MEDIA_KEY_TITLE].(string)
	ext := filepath.Ext(this.Name)
	str := CleanName(title) + ext
	if track > 0 {
		str = fmt.Sprintf("%02d - %v", track, str)
	}
	if disc > 0 {
		str = fmt.Sprintf("%02d - %v", disc, str)
	}
	return str
}

func CleanName(value string) string {
	value = strings.Replace(value, "/", "_", -1)
	value = strings.Replace(value, ".", "_", -1)
	value = strings.Replace(value, ":", "_", -1)
	return value
}
*/
