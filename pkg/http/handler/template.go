package handler

import (
	"net/http"
	"os"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Templates struct {
	gopi.Unit
	gopi.Server
	gopi.Logger
	*cache
	*renderers

	folder *string
}

type Template struct {
	gopi.Logger
	*cache
	*renderers
	name string
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Templates) Define(cfg gopi.Config) error {
	this.folder = cfg.FlagString("http.templates", "", "Path to HTML Templates")
	return nil
}

func (this *Templates) New(gopi.Config) error {
	// Where there is no templates argument
	if *this.folder == "" {
		return nil
	}

	// Read all templates
	if stat, err := os.Stat(*this.folder); os.IsNotExist(err) {
		return gopi.ErrNotFound.WithPrefix(*this.folder)
	} else if err != nil {
		return gopi.ErrBadParameter.WithPrefix(*this.folder)
	} else if stat.IsDir() == false {
		return gopi.ErrBadParameter.WithPrefix(*this.folder)
	} else if cache, err := NewCache(*this.folder); err != nil {
		return err
	} else {
		this.cache = cache
		this.renderers = NewRenderers()
	}

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// METHODS - TEMPLATES

// ServeTemplate registers a service to serve a template for a path
func (this *Templates) ServeTemplate(path, name string) error {
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("Server")
	} else if this.cache == nil {
		return gopi.ErrBadParameter.WithPrefix("-http.templates")
	} else if _, err := this.cache.Get(name); err != nil {
		return err
	} else if err := this.Server.RegisterService(path, this.NewTemplateService(name)); err != nil {
		return err
	} else {
		this.Debugf("Register Template %q => %q", path, name)
	}

	// Return success
	return nil
}

// RegisterRenderer registers a document renderer for named template
func (this *Templates) RegisterRenderer(name string, renderer gopi.HttpRenderer) error {
	if this.cache == nil {
		return gopi.ErrBadParameter.WithPrefix("-http.templates")
	} else if _, err := this.cache.Get(name); err != nil {
		return err
	} else if err := this.renderers.Register(name, renderer); err != nil {
		return err
	} else {
		this.Debugf("Register Renderer %q => %v", name, renderer)
	}

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// METHODS - TEMPLATE HANDLER

// Create a new template handler
func (this *Templates) NewTemplateService(name string) http.Handler {
	return &Template{this.Logger, this.cache, this.renderers, name}
}

// Serve a template through a renderer
func (this *Template) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var content interface{}
	var modified time.Time
	var err error

	// Get renderer
	renderer := this.renderers.Renderer(this.name)

	this.Debugf("req=%v renderer=%v", req.URL, renderer)
	this.Debugf("  template=%v", this.name)

	// Get content and modified time from cache
	content, modified = this.renderers.Get(req)
	if renderer != nil {
		if content == nil || renderer.IsModifiedSince(req, modified) {
			// Content needs to be rendered
			content, modified, err = this.renderers.Render(renderer, req)
			if err != nil {
				this.Debugf("  render content returns error: %v", err)
			} else if content == nil {
				this.Debugf("  render content returns no content")
			} else {
				this.Debugf("  render content returns modification date: %v", modified)
			}
		}

		// Deal with any errors from generating content
		if err != nil {
			this.ServeError(w, err)
			return
		} else if content == nil {
			// If content is nil then return NoContent error
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
			return
		}
	}

	// Update modification time if template modification time is later
	if modtime := this.cache.Modified(this.name); modtime.IsZero() == false {
		if modified.IsZero() || modtime.After(modified) {
			this.Debugf("  template updates modification time: %v", modtime)
			modified = modtime
		}
	}

	// Check for If-Modified-Since header
	if ifmodified := req.Header.Get("If-Modified-Since"); ifmodified != "" {
		if date, err := time.Parse(http.TimeFormat, ifmodified); err == nil {
			if date.After(modified) {
				this.Debugf("  If-Modified-Since %v: Returning %v", ifmodified, http.StatusNotModified)
				http.Error(w, http.StatusText(http.StatusNotModified), http.StatusNotModified)
				return
			}
		}
	}

	// Set cache headers
	if modified.IsZero() == false {
		w.Header().Set("Last-Modified", modified.Format(http.TimeFormat))
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}

	// Here we have content and a modified date set, so serve template
	// with content
	if tmpl, err := this.cache.Get(this.name); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err := tmpl.Execute(w, content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Serve error
func (this *Template) ServeError(w http.ResponseWriter, err error) {
	if err_, ok := err.(gopi.HttpError); ok {
		http.Error(w, err_.Error(), err_.Code())
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
