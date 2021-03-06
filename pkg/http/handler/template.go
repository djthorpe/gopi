package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	fcgi "github.com/djthorpe/gopi/v3/pkg/http/fcgi"
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

type TemplateHandler struct {
	gopi.Logger
	*TemplateCache
	*RenderCache

	path    string
	docroot string
}

type httpServer interface {
	gopi.Server

	Mux() *http.ServeMux
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
func (this *Templates) Serve(path, docroot string) error {
	// Crete a TemplateHandler
	if handler, err := this.NewTemplateHandler(path, docroot); err != nil {
		return err
	} else if err := this.Server.RegisterService(path, handler); err != nil {
		return err
	}

	// Return success
	return nil
}

// RegisterRenderer registers a document renderer
func (this *Templates) RegisterRenderer(r gopi.HttpRenderer) error {
	return this.RenderCache.Register(r)
}

// Env returns the process environment for a request
func (this *Templates) Env(req *http.Request) map[string]string {
	return fcgi.ProcessEnv(req)
}

// Render returns content and modified time for a path
func (this *Templates) Render(req *http.Request) (gopi.HttpRenderContext, error) {
	handler, _ := this.Server.(httpServer).Mux().Handler(req)
	if handler_, ok := handler.(*TemplateHandler); ok == false {
		return gopi.HttpRenderContext{}, Error(req, http.StatusNotFound)
	} else {
		return handler_.Serve(req)
	}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Templates) NewTemplateHandler(path, docroot string) (http.Handler, error) {
	h := new(TemplateHandler)
	h.TemplateCache = this.TemplateCache
	h.RenderCache = this.RenderCache
	h.Logger = this.Logger

	// Check path argument
	if strings.HasPrefix(path, "/") == false {
		return nil, gopi.ErrBadParameter.WithPrefix("NewTemplateHandler: ", path)
	} else {
		h.path = path
	}

	// Check docroot argument
	if stat, err := os.Stat(docroot); err != nil {
		return nil, gopi.ErrBadParameter.WithPrefix("NewTemplateHandler: ", err)
	} else if stat.IsDir() == false {
		return nil, gopi.ErrBadParameter.WithPrefix("NewTemplateHandler: ", docroot)
	} else {
		h.docroot = docroot
	}

	// Return success
	return h, nil
}

// Serve is internal version of renderer
func (this *TemplateHandler) Serve(req *http.Request) (gopi.HttpRenderContext, error) {
	// Get renderer or return NOT IMPLEMENTED error
	renderer := this.RenderCache.Get(this.docroot, req)
	if renderer == nil {
		return gopi.HttpRenderContext{}, Error(req, http.StatusNotImplemented)
	} else {
		return renderer.ServeContent(this.docroot, req)
	}
}

// ServeHTTP a template through a renderer
func (this *TemplateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Get renderer or return NOT IMPLEMENTED error
	renderer := this.RenderCache.Get(this.docroot, req)
	if renderer == nil {
		this.ServeError(w, Error(req, http.StatusNotImplemented))
		return
	}

	this.Debugf("ServeHTTP: req=%v renderer=%v", req.URL, renderer)

	// Check for If-Modified-Since header on content
	if ifmodified := req.Header.Get("If-Modified-Since"); ifmodified != "" {
		if date, err := time.Parse(http.TimeFormat, ifmodified); err == nil {
			if renderer.IsModifiedSince(this.docroot, req, date) == false {
				this.Debugf("  If-Modified-Since %v: Returning %v", ifmodified, http.StatusNotModified)
				this.ServeError(w, Error(req, http.StatusNotModified))
				return
			}
		}
	}

	// Render Content
	ctx, err := renderer.ServeContent(this.docroot, req)
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

// Serve error
func (this *TemplateHandler) ServeError(w http.ResponseWriter, err error) {
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

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Templates) String() string {
	str := "<http.templates"
	str += " " + this.TemplateCache.String()
	str += " " + this.RenderCache.String()
	return str + ">"
}
