package influxdb

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/djthorpe/gopi/v3"
)

type Writer struct {
	gopi.Unit

	// Parameters
	url        *string
	skipverify *bool
	timeout    *time.Duration

	// Members
	db     *url.URL
	client *http.Client
}

func (this *Writer) Define(cfg gopi.Config) error {
	this.url = cfg.FlagString("influxdb.url", "", "Database URL")
	this.skipverify = cfg.FlagBool("influxdb.slipverify", false, "Skip SSL certificate verification")
	this.timeout = cfg.FlagDuration("influxdb.timeout", 15*time.Second, "Database connection timeout")

	return nil
}

func (this *Writer) New(cfg gopi.Config) error {
	if u, err := url.Parse(*this.url); err != nil {
		return err
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return gopi.ErrBadParameter.WithPrefix("Unsupported scheme")
	} else {
		this.db = u
	}

	// Create transport
	this.client = &http.Client{
		Timeout: *this.timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: *this.skipverify,
			},
		},
	}

	// Return success
	return nil
}
