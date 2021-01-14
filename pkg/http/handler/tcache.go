package handler

import (
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type TemplateCache struct {
	sync.RWMutex
	gopi.Unit
	gopi.Logger

	folder *string
	t      map[string]tcached
}

type tcached struct {
	os.FileInfo
	*template.Template
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *TemplateCache) Define(cfg gopi.Config) error {
	this.folder = cfg.FlagString("http.templates", "", "Path to HTML Templates")
	return nil
}

func (this *TemplateCache) New(gopi.Config) error {
	this.Require(this.Logger)

	// Where there is no templates argument provided, the unit does
	// not activate
	if *this.folder == "" {
		return nil
	}

	// Check folder argument
	if stat, err := os.Stat(*this.folder); os.IsNotExist(err) {
		return gopi.ErrNotFound.WithPrefix(*this.folder)
	} else if err != nil {
		return gopi.ErrBadParameter.WithPrefix(*this.folder)
	} else if stat.IsDir() == false {
		return gopi.ErrBadParameter.WithPrefix(*this.folder)
	}

	// Read all templates in from a folder, every template needs to
	// be parsed without error
	files, err := ioutil.ReadDir(*this.folder)
	if err != nil {
		return err
	}

	// Create cache of templates
	this.t = make(map[string]tcached, len(files))

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
		path := filepath.Join(*this.folder, file.Name())
		if tmpl, err := template.New(file.Name()).Funcs(this.funcmap()).ParseFiles(path); err != nil {
			result = multierror.Append(result, err)
		} else {
			tmpl = tmpl.Funcs(this.funcmap())
			key := tmpl.Name()
			this.t[key] = tcached{file, tmpl}
			this.Debugf("Parsed template: %q", tmpl.Name())
		}
	}

	// Return any errors
	return result
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Get returns a template keyed by name. If the template has been
// updated on the filesystem, then it is reparsed.
func (this *TemplateCache) Get(name string) (*template.Template, time.Time, error) {
	t := this.get(name)
	path := filepath.Join(*this.folder, name)

	// No template exists with this name, return nil
	if t.Template == nil {
		return nil, time.Time{}, gopi.ErrNotFound.WithPrefix(name)
	}

	// Check case where template has not changed
	info, err := os.Stat(path)
	if err != nil {
		return t.Template, time.Time{}, err
	} else if info.ModTime() == t.ModTime() {
		return t.Template, t.ModTime(), nil
	}

	// Re-parse template and return it
	tmpl, err := template.New(name).Funcs(this.funcmap()).ParseFiles(path)
	if err == nil {
		this.Debugf("Reparsed template: %q", tmpl.Name())
		this.set(name, tmpl, info)
		return tmpl, info.ModTime(), nil
	} else {
		this.Debugf("Parse Error: %q: %v", name, err)
		return t.Template, t.ModTime(), err
	}
}

// Modified returns the modification time of a template by name
// or returns zero-time if the template does not exist
func (this *TemplateCache) Modified(name string) time.Time {
	t := this.get(name)
	if t.Template == nil {
		return time.Time{}
	} else {
		return t.ModTime()
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *TemplateCache) String() string {
	str := "<http.templatecache"
	if *this.folder != "" {
		str += fmt.Sprintf(" folder=%q", *this.folder)
	}
	for k := range this.t {
		str += " tmpl=" + strconv.Quote(k)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *TemplateCache) get(name string) tcached {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.t == nil {
		return tcached{}
	} else if t, exists := this.t[name]; exists == false {
		return tcached{}
	} else {
		return t
	}
}

func (this *TemplateCache) set(name string, t *template.Template, info os.FileInfo) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	this.t[name] = tcached{info, t}
}

func (this *TemplateCache) funcmap() template.FuncMap {
	return template.FuncMap{
		"ssi":         funcSSI,
		"pathescape":  funcPathEscape,
		"queryescape": funcQueryEscape,
		"textescape":  funcTextEscape,
	}
}

func funcPathEscape(value string) string {
	return url.PathEscape(value)
}

func funcQueryEscape(value string) string {
	return url.QueryEscape(value)
}

func funcTextEscape(value string) string {
	return html.EscapeString(value)
}

func funcSSI(cmd string, args ...string) template.HTML {
	switch cmd {
	case "block":
		if len(args) == 1 {
			return template.HTML(fmt.Sprintf("<!--# block name=%q -->", args[0]))
		}
	case "endblock", "else", "endif":
		if len(args) == 0 {
			return template.HTML(fmt.Sprintf("<!--# %v -->", cmd))
		}
	case "echo":
		if len(args) == 1 {
			return template.HTML(fmt.Sprintf("<!--# var name=%q -->", args[0]))
		} else if len(args) == 2 {
			return template.HTML(fmt.Sprintf("<!--# var name=%q default=%q -->", args[0], args[1]))
		}
	case "if", "elif":
		if len(args) == 1 {
			return template.HTML(fmt.Sprintf("<!--# %v expr=%q -->", cmd, args[0]))
		}
	case "include":
		if len(args) == 1 {
			return template.HTML(fmt.Sprintf("<!--# %v virtual=%q -->", cmd, args[0]))
		}
		if len(args) == 2 && args[1] == "wait" {
			return template.HTML(fmt.Sprintf("<!--# %v virtual=%q wait=%q -->", cmd, args[0], "yes"))
		}
	case "set":
		if len(args) == 1 {
			return template.HTML(fmt.Sprintf("<!--# set var=%q value=%q -->", args[0], ""))
		}
		if len(args) == 2 {
			return template.HTML(fmt.Sprintf("<!--# set var=%q value=%q -->", args[0], args[1]))
		}
	}
	return template.HTML(fmt.Sprintf("[an error occurred while processing the directive: %q ]", cmd))
}
