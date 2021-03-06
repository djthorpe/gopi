package handler

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	gopi "github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Static struct {
	gopi.Unit
	gopi.Server
	gopi.Logger

	folder *string
}

type static struct {
	http.Handler
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Static) Define(cfg gopi.Config) error {
	this.folder = cfg.FlagString("http.static", "", "Path to static folder")
	return nil
}

func (this *Static) New(gopi.Config) error {
	// Where there is no static argument, use current working
	// directory
	if *this.folder == "" {
		if wd, err := os.Getwd(); err != nil {
			return err
		} else {
			*this.folder = wd
		}
	}

	// Check static folder
	if stat, err := os.Stat(*this.folder); os.IsNotExist(err) {
		return gopi.ErrNotFound.WithPrefix(*this.folder)
	} else if err != nil {
		return gopi.ErrBadParameter.WithPrefix(*this.folder)
	} else if stat.IsDir() == false {
		return gopi.ErrBadParameter.WithPrefix(*this.folder)
	}

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// METHODS

// Serve registers a service to serve static files from all folders
// under the named path
func (this *Static) Serve(path string) error {
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("Serve")
	} else if *this.folder == "" {
		return gopi.ErrBadParameter.WithPrefix("Serve")
	}

	files, err := ioutil.ReadDir(*this.folder)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() == false {
			continue
		}
		child := filepath.Join(path, file.Name())
		if strings.HasSuffix(child, "/") == false {
			child = child + "/"
		}
		folder := filepath.Join(*this.folder, child)
		if err := this.Server.RegisterService(child, NewStaticHandler(child, folder)); err != nil {
			return err
		} else {
			this.Debugf("Register Static %q => %v", child, folder)
		}
	}

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// HANDLER

func NewStaticHandler(root, folder string) http.Handler {
	this := new(static)
	this.Handler = http.StripPrefix(root, http.FileServer(http.Dir(folder)))
	return this
}

func (this *static) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: Prevent listing of directories
	this.Handler.ServeHTTP(w, req)
}
