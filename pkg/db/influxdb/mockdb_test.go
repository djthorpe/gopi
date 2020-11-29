package influxdb_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
)

type server struct {
	sync.WaitGroup
	*http.Server
	*testing.T

	errs chan error
}

func NewMockServer(t *testing.T, addr string) (*server, error) {
	this := new(server)
	mux := http.NewServeMux()
	this.Server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	this.errs = make(chan error)
	this.T = t

	// Setup handlers
	mux.HandleFunc("/ping", this.Ping)
	mux.HandleFunc("/write", this.Write)

	// Start serving in the background
	this.WaitGroup.Add(1)
	go func() {
		this.errs <- this.Server.ListenAndServe()
		this.WaitGroup.Done()
	}()

	// Return
	return this, nil
}

func (this *server) Close() error {
	var result error

	// Shutdown and wait for end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := this.Server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		result = multierror.Append(result, err)
	}

	// Receive shutdown information
	if err := <-this.errs; err != nil && err != http.ErrServerClosed {
		result = multierror.Append(result, err)
	}

	// Wait for end of goroutine
	this.WaitGroup.Wait()

	// Close error channel
	close(this.errs)

	// Return any errors
	return result
}

func (this *server) Ping(w http.ResponseWriter, r *http.Request) {
	this.T.Log("Mock server got Ping")
	http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
}

func (this *server) Write(w http.ResponseWriter, r *http.Request) {
	this.T.Log("Mock server got Write")
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		this.T.Error(err)
	} else {
		this.T.Log(string(data))
	}
	r.Body.Close()
	http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
}
