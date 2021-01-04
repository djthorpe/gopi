package handler

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type cache struct {
	sync.RWMutex
	folder string
	t      map[string]cached
}

type cached struct {
	os.FileInfo
	*template.Template
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCache(folder string) (*cache, error) {
	// Read files in folder
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	// Create cache
	this := new(cache)
	this.t = make(map[string]cached, len(files))
	this.folder = folder

	// Read templates into cache. Collect parse errors
	// into results
	var result error
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.Mode().IsRegular() == false {
			continue
		}
		if tmpl, err := template.ParseFiles(filepath.Join(folder, file.Name())); err != nil {
			result = multierror.Append(result, err)
		} else {
			key := tmpl.Name()
			this.t[key] = cached{file, tmpl}
		}
	}

	// Return cache or error
	if result != nil {
		return nil, result
	} else {
		return this, nil
	}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *cache) Get(name string) (*template.Template, error) {
	t := this.get(name)
	path := filepath.Join(this.folder, name)

	// No template exists with this name, return nil
	if t.Template == nil {
		return nil, gopi.ErrNotFound.WithPrefix(name)
	}

	// Check case where template has not changed
	info, err := os.Stat(path)
	if err != nil {
		return t.Template, err
	} else if info.ModTime() == t.ModTime() {
		return t.Template, nil
	}

	// Re-parse template and return it
	tmpl, err := template.ParseFiles(path)
	if err == nil {
		return this.set(name, tmpl, info), nil
	} else {
		return t.Template, err
	}
}

func (this *cache) Modified(name string) time.Time {
	t := this.get(name)
	if t.Template == nil {
		return time.Time{}
	} else {
		return t.ModTime()
	}
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *cache) get(name string) cached {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.t == nil {
		return cached{}
	} else if t, exists := this.t[name]; exists == false {
		return cached{}
	} else {
		return t
	}
}

func (this *cache) set(name string, t *template.Template, info os.FileInfo) *template.Template {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	this.t[name] = cached{info, t}
	return t
}
