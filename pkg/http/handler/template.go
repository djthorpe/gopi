package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Templates struct {
	gopi.Unit
	gopi.Server
	gopi.Logger
	*TemplateCache
	*RenderCache
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Templates) New(gopi.Config) error {
	this.Require(this.Server, this.Logger, this.TemplateCache, this.RenderCache)

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Serve registers a service to serve templates for a path
func (this *Templates) Serve(path string) error {
	if err := this.Server.RegisterService(path, this); err != nil {
		return err
	}

	// Return success
	return nil
}

// RegisterRenderer registers a document renderer
func (this *Templates) RegisterRenderer(r gopi.HttpRenderer) error {
	return this.RenderCache.Register(r)
}

// Serve error
func (this *Templates) ServeError(w http.ResponseWriter, err error) {
	if err_, ok := err.(gopi.HttpError); ok == false {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if err_.Code() == http.StatusPermanentRedirect || err_.Code() == http.StatusTemporaryRedirect {
		this.Debugf("  Code: %v", err_.Error())
		this.Debugf("  Location: %v", err_.Path())

		w.Header().Set("Location", err_.Path())
		http.Error(w, err_.Error(), err_.Code())
	} else {
		http.Error(w, err_.Error(), err_.Code())
	}
}

// Serve a template through a renderer
func (this *Templates) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Get renderer or return NOT IMPLEMENTED error
	renderer := this.RenderCache.Get(req)
	if renderer == nil {
		this.ServeError(w, Error(req, http.StatusNotImplemented))
		return
	}

	this.Debugf("ServeHTTP: req=%v renderer=%v", req.URL, renderer)

	// Check for If-Modified-Since header on content
	if ifmodified := req.Header.Get("If-Modified-Since"); ifmodified != "" {
		if date, err := time.Parse(http.TimeFormat, ifmodified); err == nil {
			if renderer.IsModifiedSince(req, date) == false {
				this.Debugf("  If-Modified-Since %v: Returning %v", ifmodified, http.StatusNotModified)
				this.ServeError(w, Error(req, http.StatusNotModified))
				return
			}
		}
	}

	// Render Content
	ctx, err := renderer.ServeContent(req)
	if err != nil {
		this.ServeError(w, err)
		return
	} else if ctx.Content == nil {
		this.ServeError(w, Error(req, http.StatusNoContent))
		return
	}

	// Get template and template modification time
	var tmpl *template.Template
	var modtime time.Time
	if ctx.Template != "" {
		if tmpl, modtime, err = this.TemplateCache.Get(ctx.Template); err != nil {
			this.Debugf("  Template %q: Error: %v", ctx.Template, err)
			this.ServeError(w, Error(req, http.StatusNotFound, err.Error()))
			return
		}
	}

	// Update modification time for page if the template is later
	if ctx.Modified.IsZero() == false && modtime.After(ctx.Modified) {
		this.Debugf("  Template updates modification time: %v", modtime)
		ctx.Modified = modtime
	}

	// Set default type
	if ctx.Type == "" {
		ctx.Type = "application/octet-stream"
	}

	// Set headers
	w.Header().Set("Content-Type", ctx.Type)
	if ctx.Modified.IsZero() == false {
		w.Header().Set("Last-Modified", ctx.Modified.Format(http.TimeFormat))
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}

	// If no template then we expect the content to be []byte
	if tmpl == nil {
		if data, ok := ctx.Content.([]byte); ok {
			w.Header().Set("Content-Length", fmt.Sprint(len(data)))
			w.WriteHeader(http.StatusOK)
			if req.Method != http.MethodHead {
				w.Write(data)
			}
			return
		} else {
			this.ServeError(w, Error(req, http.StatusInternalServerError))
			return
		}
	}

	// Debugging
	/*
		if this.Logger.IsDebug() {
			if json, err := json.MarshalIndent(ctx.Content, "  ", "  "); err == nil {
				this.Debugf(string(json))
			}
		}
	*/

	// Execute through a template
	data := new(bytes.Buffer)
	if err := tmpl.Execute(data, ctx.Content); err != nil {
		this.ServeError(w, Error(req, http.StatusInternalServerError, err.Error()))
		return
	}

	// Set content length and write data
	w.Header().Set("Content-Length", fmt.Sprint(data.Len()))
	w.WriteHeader(http.StatusOK)
	if req.Method != http.MethodHead {
		w.Write(data.Bytes())
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Templates) String() string {
	str := "<http.templates"
	str += " " + this.TemplateCache.String()
	str += " " + this.RenderCache.String()
	return str + ">"
}
