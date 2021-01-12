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
)

/////////////////////////////////////////////////////////////////////
// TYPES

type HttpIndexRenderer struct {
	folder   string
	template string
}

type HttpIndexContent struct {
	Path    string         `json:"path"`
	Content []IndexContent `json:"content"`
}

type IndexContent struct {
	Name    string    `json:"name"`
	IsDir   bool      `json:"isdir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modtime"`
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewIndexRenderer(folder, template string) gopi.HttpRenderer {
	this := new(HttpIndexRenderer)
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

	// Return success
	return this
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *HttpIndexRenderer) ServeContent(req *http.Request) (gopi.HttpRenderContext, error) {
	path := filepath.Join(this.folder, req.URL.Path)

	// Update modified time based on all eligible files
	content := HttpIndexContent{
		Path: req.URL.Path,
	}

	if mtime, err := filesForFolder(path, func(file os.FileInfo) error {
		content.Content = append(content.Content, IndexContent{
			Name:    file.Name(),
			IsDir:   file.IsDir(),
			Size:    file.Size(),
			ModTime: file.ModTime(),
		})
		return nil
	}); err != nil {
		return gopi.HttpRenderContext{}, err
	} else {
		return gopi.HttpRenderContext{
			Template: this.template,
			Type:     mime.TypeByExtension(".html"),
			Content:  content,
			Modified: mtime,
		}, nil
	}
}

func (this *HttpIndexRenderer) IsModifiedSince(req *http.Request, t time.Time) bool {
	path := filepath.Join(this.folder, req.URL.Path)

	// Update modified time based on all eligible files
	if mtime, err := filesForFolder(path, nil); err != nil {
		return false
	} else {
		return mtime.After(t)
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *HttpIndexRenderer) String() string {
	str := "<http.indexrenderer"
	if this.folder != "" {
		str += fmt.Sprintf(" folder=%q", this.folder)
	}
	if this.template != "" {
		str += fmt.Sprintf(" template=%q", this.template)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func filesForFolder(path string, fn func(os.FileInfo) error) (time.Time, error) {
	var mtime time.Time

	// Check path
	if stat, err := os.Stat(path); err != nil {
		return time.Time{}, err
	} else if stat.IsDir() == false {
		return time.Time{}, gopi.ErrBadParameter
	} else {
		mtime = stat.ModTime()
	}

	// Walk through files in this folder in alphabetical order
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return time.Time{}, err
	}

	// Render files, ignoring any subsequent folders
	for _, file := range files {
		if file.Mode().IsRegular() == false && file.Mode().IsDir() == false {
			continue
		}
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.ModTime().After(mtime) {
			mtime = file.ModTime()
		}
		if fn != nil {
			if err := fn(file); err != nil {
				return time.Time{}, err
			}
		}
	}

	// Return success
	return mtime, nil
}
