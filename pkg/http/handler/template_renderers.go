package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type renderers struct {
	sync.RWMutex

	r map[string]gopi.HttpRenderer
	c map[string]content
}

type content struct {
	modified time.Time
	content  interface{}
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewRenderers() *renderers {
	this := new(renderers)
	this.r = make(map[string]gopi.HttpRenderer)
	this.c = make(map[string]content)
	return this
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *renderers) Register(name string, renderer gopi.HttpRenderer) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if renderer == nil {
		return gopi.ErrBadParameter.WithPrefix("Register")
	} else if _, exists := this.r[name]; exists {
		return gopi.ErrDuplicateEntry.WithPrefix("Register", name)
	} else {
		this.r[name] = renderer
	}

	// Return success
	return nil
}

func (this *renderers) Get(req *http.Request) (interface{}, time.Time) {
	if req == nil || req.URL == nil {
		return nil, time.Time{}
	} else {
		return this.get(keyForRequest(req))
	}
}

func (this *renderers) Set(req *http.Request, content interface{}, modified time.Time) {
	if req != nil && req.URL != nil {
		this.set(keyForRequest(req), content, modified)
	}
}

func (this *renderers) Renderer(name string) gopi.HttpRenderer {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if r, exists := this.r[name]; exists {
		return r
	} else {
		return nil
	}
}

func (this *renderers) Render(renderer gopi.HttpRenderer, req *http.Request) (interface{}, time.Time, error) {
	// Check renderer
	if renderer == nil {
		return nil, time.Time{}, gopi.ErrNotFound
	}
	// Get existing content and modified date
	key := keyForRequest(req)
	if content, modified := this.get(key); content != nil {
		// Serve from cache if not modified
		if renderer.IsModifiedSince(modified) == false {
			return content, modified, nil
		}
	}

	// Generate content
	if content, modified, err := renderer.ServeContent(req); err != nil {
		// If error then remove content from cache
		this.set(key, nil, time.Time{})

		// Return the error
		return nil, time.Time{}, err
	} else {
		if content == nil {
			// Delete content from cache if content is nil
			this.set(key, nil, time.Time{})
		} else if modified.IsZero() == false {
			this.set(key, content, modified)
		}
		// Return content
		return content, modified, nil
	}
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *renderers) get(key string) (interface{}, time.Time) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if cached, exists := this.c[key]; exists {
		return cached.content, cached.modified
	} else {
		return nil, time.Time{}
	}
}

func (this *renderers) set(key string, value interface{}, modified time.Time) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	if value != nil && modified.IsZero() == false {
		this.c[key] = content{modified, value}
	} else {
		delete(this.c, key)
	}
}

func keyForRequest(req *http.Request) string {
	return req.URL.String()
}
