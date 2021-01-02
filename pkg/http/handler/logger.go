package handler

import (
	"net"
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
		this.Metrics.Field("url", ""),
		this.Metrics.Field("host", ""),
		this.Metrics.Field("path", ""),
		this.Metrics.Field("query", ""),
		this.Metrics.Field("remoteAddr", ""),
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
	tags := []gopi.Field{
		this.Metrics.Field("method", req.Method),
		this.Metrics.Field("url", req.RequestURI),
		this.Metrics.Field("host", req.Host),
		this.Metrics.Field("path", req.URL.Path),
		this.Metrics.Field("query", req.URL.RawQuery),
	}
	if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		tags = append(tags, this.Metrics.Field("remoteAddr", host))
	}
	/*	for k, v := range w.Header() {
			fmt.Println("response header", k, v)
		}
	*/
	// Emit Metrics
	this.Metrics.EmitTS(this.name, now, tags, time.Since(now).Seconds(), req.UserAgent())
}
