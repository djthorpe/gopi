package renderer

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/http/handler"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type HttpTextRenderer struct {
	ext      map[string]bool
	folder   string
	template string
}

type TextContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EXT = ".txt .tmpl .md .go"
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewTextRenderer(folder, template string) gopi.HttpRenderer {
	this := new(HttpTextRenderer)
	if folder == "" {
		return nil
	} else if stat, err := os.Stat(folder); err != nil {
		return nil
	} else if stat.IsDir() == false {
		return nil
	} else {
		this.folder = folder
	}
	if template == "" {
		return nil
	} else {
		this.template = template
	}

	// Set extension map
	this.ext = make(map[string]bool, 10)
	for _, ext := range strings.Fields(EXT) {
		this.ext[ext] = true
	}

	// Return success
	return this
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *HttpTextRenderer) IsModifiedSince(req *http.Request, t time.Time) bool {
	// Check file extension
	ext := filepath.Ext(req.URL.Path)
	if _, exists := this.ext[ext]; exists == false {
		return false
	}
	// Must exist and be a regular file, not starting with a period "."
	path := filepath.Join(this.folder, req.URL.Path)
	if stat, err := os.Stat(path); err != nil {
		return false
	} else if stat.Mode().IsRegular() == false {
		return false
	} else if strings.HasPrefix(stat.Name(), ".") {
		return false
	} else if t.IsZero() {
		return true
	} else if t.After(stat.ModTime()) {
		return false
	} else {
		return true
	}
}

func (this *HttpTextRenderer) ServeContent(req *http.Request) (gopi.HttpRenderContext, error) {
	path := filepath.Join(this.folder, req.URL.Path)

	// Check file
	stat, err := os.Stat(path)
	if err != nil {
		return gopi.HttpRenderContext{}, handler.Error(req, http.StatusNotFound, err.Error())
	} else if stat.Mode().IsRegular() == false {
		return gopi.HttpRenderContext{}, handler.Error(req, http.StatusNotFound, err.Error())
	} else if strings.HasPrefix(stat.Name(), ".") {
		return gopi.HttpRenderContext{}, handler.Error(req, http.StatusForbidden, err.Error())
	}

	// Read file contents
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return gopi.HttpRenderContext{}, handler.Error(req, http.StatusInternalServerError, err.Error())
	}

	// If query is ?static return data with no template
	// or else return file with template
	switch req.URL.RawQuery {
	case "static":
		return gopi.HttpRenderContext{
			Type:     mime.TypeByExtension(".txt"),
			Content:  data,
			Modified: stat.ModTime(),
		}, nil
	default:
		return gopi.HttpRenderContext{
			Template: this.template,
			Type:     mime.TypeByExtension(".html"),
			Content:  TextContent{req.URL.Path, string(data)},
			Modified: stat.ModTime(),
		}, nil
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *HttpTextRenderer) String() string {
	str := "<http.textrenderer"
	if this.folder != "" {
		str += fmt.Sprintf(" folder=%q", this.folder)
	}
	if this.template != "" {
		str += fmt.Sprintf(" template=%q", this.template)
	}
	return str + ">"
}
