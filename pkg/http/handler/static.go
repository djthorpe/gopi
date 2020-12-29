package handler

import (
	"net/http"

	gopi "github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Static struct {
	gopi.Unit
	gopi.Server
}

// Register a service to serve static files to path from a folder
func (this *Static) ServeFolder(path, folder string) error {
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("ServeFolder")
	} else if err := this.Server.RegisterService(path, http.FileServer(http.Dir(folder))); err != nil {
		return err
	}

	// Return success
	return nil
}
