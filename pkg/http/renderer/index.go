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

type HttpIndexRenderer struct {
	gopi.Unit
	gopi.Logger
	gopi.HttpTemplate

	template string
}

type HttpIndexContent struct {
	Path    string         `json:"path"`
	Content []IndexContent `json:"content"`
}

type IndexContent struct {
	Path    string    `json:"path"`
	Name    string    `json:"name"`
	IsDir   bool      `json:"isdir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modtime"`
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *HttpIndexRenderer) New(gopi.Config) error {
	this.Require(this.Logger, this.HttpTemplate)

	if err := this.HttpTemplate.RegisterRenderer(this); err != nil {
		return err
	} else {
		this.template = "index.tmpl"
	}

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *HttpIndexRenderer) ServeContent(docroot string, req *http.Request) (gopi.HttpRenderContext, error) {
	// Add a slash if not at the end
	if strings.HasSuffix(req.URL.Path, "/") == false {
		return gopi.HttpRenderContext{}, handler.Redirect(req, http.StatusPermanentRedirect, req.URL.Path+"/")
	}

	// Compute physical path
	path := filepath.Join(docroot, req.URL.Path)

	// Update modified time based on all eligible files
	content := HttpIndexContent{
		Path: req.URL.Path,
	}

	if mtime, err := filesForFolder(path, func(file os.FileInfo) error {
		content.Content = append(content.Content, IndexContent{
			Path:    filepath.Join(content.Path, file.Name()),
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

func (this *HttpIndexRenderer) IsModifiedSince(docroot string, req *http.Request, t time.Time) bool {
	path := filepath.Join(docroot, req.URL.Path)

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
