package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	gopi "github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Templates struct {
	gopi.Unit
	gopi.Server

	folder *string
}

type Template struct {
	Name string
}

func (this *Templates) Define(cfg gopi.Config) error {
	this.folder = cfg.FlagString("http.templates", "", "Path to HTML Templates")
	return nil
}

func (this *Templates) New(gopi.Config) error {
	if *this.folder != "" {
		if stat, err := os.Stat(*this.folder); os.IsNotExist(err) {
			return gopi.ErrNotFound.WithPrefix(*this.folder)
		} else if err != nil {
			return gopi.ErrBadParameter.WithPrefix(*this.folder)
		} else if stat.IsDir() == false {
			return gopi.ErrBadParameter.WithPrefix(*this.folder)
		} else if files, err := ioutil.ReadDir(*this.folder); err != nil {
			return err
		} else {
			for _, file := range files {
				if strings.HasPrefix(file.Name(), ".") {
					continue
				}
				fmt.Println("TODO: Template: ", file)
			}
		}
	}

	// Return success
	return nil
}

// Register a service to serve a template for a path
func (this *Templates) ServeTemplate(path, template string) error {
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("Server")
	} else if err := this.Server.RegisterService(path, this.NewTemplateService(template)); err != nil {
		return err
	}

	// Return success
	return nil
}

// Create a new template handler
func (this *Templates) NewTemplateService(name string) http.Handler {
	return &Template{name}
}

func (this *Template) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO
}
