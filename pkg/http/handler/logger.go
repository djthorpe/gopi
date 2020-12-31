package handler

import (
	"net/http"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Logger struct {
	gopi.Unit
	gopi.Server
	gopi.Metrics
}

type LoggerHandler struct {
	*Logger
	name string
	http.Handler
}

// Register a service which logs metrics
func (this *Logger) Log(name string) error {
	tags := []gopi.Field{
		this.Metrics.Field("method", ""),
		this.Metrics.Field("uri", ""),
		this.Metrics.Field("host", ""),
		this.Metrics.Field("path", ""),
		this.Metrics.Field("query", ""),
		this.Metrics.Field("remoteHost", ""),
	}
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("Server")
	} else if this.Metrics == nil {
		return gopi.ErrInternalAppError.WithPrefix("Metrics")
	} else if _, err := this.Metrics.NewMeasurement(name, "latency float64,useragent string", tags...); err != nil {
		return err
	} else if err := this.Server.RegisterService(nil, &LoggerHandler{this, name, nil}); err != nil {
		return err
	}

	// Return success
	return nil
}

// Set handler
func (this *LoggerHandler) SetHandler(handler http.Handler) {
	this.Handler = handler
}

// Log metrics
func (this *LoggerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	now := time.Now()
	if this.Handler != nil {
		this.Handler.ServeHTTP(w, req)
	}

	// TAGS
	/*
		fmt.Println("method=", req.Method)
		fmt.Println("uri=", req.RequestURI)
		fmt.Println("host=", req.Host)
		fmt.Println("path=", req.URL.Path)
		fmt.Println("query=", req.URL.RawQuery)
		if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			fmt.Println("remoteHost=", host)
		}
	*/
	/*	for k, v := range w.Header() {
			fmt.Println("response header", k, v)
		}
	*/
	// Emit Metrics
	this.Metrics.EmitTS(this.name, now, time.Since(now).Seconds(), req.UserAgent())
}
