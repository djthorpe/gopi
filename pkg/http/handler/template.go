package handler

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Templates struct {
	gopi.Unit
	gopi.Server
	gopi.Logger
	*cache

	folder *string
}

type cached struct {
	os.FileInfo
	*template.Template
}

type cache struct {
	sync.RWMutex

	folder string
	t      map[string]cached
}

type Template struct {
	*cache
	Name string
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
	}

	// Return success
	return nil
}

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

	// Read templates into cache
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
	if result != nil {
		return nil, result
	} else {
		return this, nil
	}
}

/////////////////////////////////////////////////////////////////////
// METHODS - CACHE

func (this *cache) Get(name string) (*template.Template, error) {
	t := this.get(name)
	path := filepath.Join(this.folder, name)
	if t.Template == nil {
		return nil, gopi.ErrNotFound.WithPrefix(name)
	}

	info, err := os.Stat(path)
	if err != nil {
		return t.Template, err
	} else if info.ModTime() == t.ModTime() {
		return t.Template, nil
	}

	tmpl, err := template.ParseFiles(path)
	if err == nil {
		return this.set(name, tmpl, info), nil
	} else {
		return t.Template, err
	}
}

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

/////////////////////////////////////////////////////////////////////
// METHODS - TEMPLATES

// Register a service to serve a template for a path
func (this *Templates) ServeTemplate(path, name string) error {
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("Server")
	} else if this.cache == nil {
		return gopi.ErrBadParameter.WithPrefix("-http.templates")
	} else if _, err := this.cache.Get(name); err != nil {
		return err
	} else if err := this.Server.RegisterService(path, this.NewTemplateService(name)); err != nil {
		return err
	}

	// Return success
	return nil
}

// Create a new template handler
func (this *Templates) NewTemplateService(name string) http.Handler {
	return &Template{this.cache, name}
}

func (this *Template) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if tmpl, err := this.cache.Get(this.Name); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
