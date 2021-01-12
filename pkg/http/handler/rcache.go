package handler

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type RenderCache struct {
	gopi.Unit
	gopi.Logger
	sync.RWMutex
	sync.WaitGroup

	r []gopi.HttpRenderer
	c map[string]gopi.HttpRenderer
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Maximum number of entries to store for map between request
	// to renderer
	rcacheMaxSize = 1000
	rcacheMinSize = rcacheMaxSize - 100
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *RenderCache) New(gopi.Config) error {

	// Set up req->render map
	this.c = make(map[string]gopi.HttpRenderer, rcacheMaxSize)

	// Return success
	return nil
}

// Dispose resources related to the render cache
func (this *RenderCache) Dispose() error {
	// Wait for any goroutines to complete
	this.WaitGroup.Wait()

	// Exclusive lock
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Dispose resources
	this.c = nil
	this.r = nil

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Register appends a renderer to the set of renderers
func (this *RenderCache) Register(renderer gopi.HttpRenderer) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if renderer == nil {
		return gopi.ErrBadParameter.WithPrefix("Register")
	} else {
		this.r = append(this.r, renderer)
	}

	// Return success
	return nil
}

// Get returns a renderer for an incoming request, or nil if there is
// no renderer which matches the request
func (this *RenderCache) Get(req *http.Request) gopi.HttpRenderer {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Check incoming parameters
	if req == nil || req.URL == nil {
		return nil
	}

	// Check in the existing cache
	key := keyForRequest(req)
	if r, exists := this.c[key]; exists {
		return r
	}

	// Round robin the renderers to find the correct one
	for _, renderer := range this.r {
		if renderer.IsModifiedSince(req, time.Time{}) {
			// If we have found the right renderer, then cache it
			// for this specific request. Also prunes the size of the map
			// so we do it in the background
			this.WaitGroup.Add(1)
			go func() {
				this.setr(key, renderer)
				this.WaitGroup.Done()
			}()
			return renderer
		}
	}

	// No renderer found so return nil
	return nil
}

// Render will return the render context for a renderer and a request
func (this *RenderCache) Render(renderer gopi.HttpRenderer, req *http.Request) (gopi.HttpRenderContext, error) {
	// Check incoming arguments
	if renderer == nil || req == nil {
		return gopi.HttpRenderContext{}, gopi.ErrNotFound
	}

	// Get existing content and modified date
	key := keyForRequest(req)
	if ctx := this.getc(key); ctx.Content != nil {
		// Serve from cache if not modified
		if renderer.IsModifiedSince(req, ctx.Modified) == false {
			return ctx, nil
		}
	}

	// Generate content
	if ctx, err := renderer.ServeContent(req); err != nil {
		// Remove the content from the cache
		this.setc(key, gopi.HttpRenderContext{})
		// Return the error
		return gopi.HttpRenderContext{}, err
	} else {
		// Store content if context is cachable
		if ctx.Content != nil && ctx.Modified.IsZero() == false {
			this.setc(key, ctx)
		}
		// Return the content
		return ctx, nil
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *RenderCache) String() string {
	str := "<http.rendercache"
	for _, r := range this.r {
		str += " r=" + fmt.Sprint(r)
	}
	if sz := len(this.c); sz > 0 {
		str += fmt.Sprint(" cachesize=", sz)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Get returns cached content or nil
func (this *RenderCache) getc(key string) gopi.HttpRenderContext {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return gopi.HttpRenderContext{}
}

// Set saves content for a key
func (this *RenderCache) setc(string, gopi.HttpRenderContext) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	// NOOP
}

// Setr saves a mapping between request key and renderer, and prunes the
// size of the map
func (this *RenderCache) setr(key string, renderer gopi.HttpRenderer) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Prune
	if len(this.c) >= rcacheMaxSize {
		for k := range this.c {
			if len(this.c) < rcacheMinSize {
				// End loop when the size of the map is below maximum size
				break
			}
			// We assume that the keys returned are in a random order?
			this.Debugf("Pruning RenderCache for %q", k)
			delete(this.c, k)
		}
	}

	// Set
	if renderer != nil {
		this.c[key] = renderer
	}
}

func keyForRequest(req *http.Request) string {
	return req.URL.String()
}
