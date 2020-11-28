package influxdb

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Writer struct {
	gopi.Unit
	gopi.Logger
	sync.Mutex

	// Parameters
	url        *string
	skipverify *bool
	timeout    *time.Duration

	// Members
	*http.Client
	endpoint
	version string
}

type endpoint struct {
	url.URL
	db, user, password string
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultScheme   = "http"
	DefaultEndpoint = DefaultScheme + "://localhost/metrics"
	DefaultPort     = 8086
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *Writer) Define(cfg gopi.Config) error {
	this.url = cfg.FlagString("influxdb.url", "", "Database URL")
	this.skipverify = cfg.FlagBool("influxdb.slipverify", false, "Skip SSL certificate verification")
	this.timeout = cfg.FlagDuration("influxdb.timeout", 15*time.Second, "Database connection timeout")
	return nil
}

func (this *Writer) New(cfg gopi.Config) error {
	// Check URL parameter
	if endpoint, err := parseUrl(*this.url); err != nil {
		return err
	} else {
		this.endpoint = endpoint
	}

	// Create transport
	this.Client = &http.Client{
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

func (this *Writer) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Close connections
	this.Client.CloseIdleConnections()

	// Release resources
	this.Client = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Writer) Endpoint() *url.URL {
	return &this.endpoint.URL
}

func (this *Writer) Database() string {
	return this.endpoint.db
}

func (this *Writer) Credentials() (string, string) {
	return this.endpoint.user, this.endpoint.password
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Writer) Ping() (time.Duration, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Set up request
	now := time.Now()
	ep := this.endpoint.URL
	ep.Path = "/ping"
	req, err := http.NewRequest(http.MethodGet, ep.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", err, ep.String())
	}

	// Set credentials
	if this.endpoint.user != "" {
		req.SetBasicAuth(this.endpoint.user, this.endpoint.password)
	}

	// Perform the request
	resp, err := this.Client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", err, ep.String())
	}
	defer resp.Body.Close()

	// Read the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Check status code
	if resp.StatusCode != http.StatusNoContent {
		return 0, fmt.Errorf("%v: %q", strings.TrimSpace(string(body)), ep.String())
	}

	// Set version
	if version := resp.Header.Get("X-Influxdb-Version"); version != "" {
		this.Debug("X-Influxdb-Version: ", version)
		this.version = version
	}

	// Return success
	return time.Since(now), nil
}

// Write measurements to the endpoint
func (this *Writer) Write(m ...gopi.Measurement) error {
	// Perform a ping if not already done
	if this.version == "" {
		if _, err := this.Ping(); err != nil {
			return err
		}
	}

	// Lock for write
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Add measurements
	buffer := new(bytes.Buffer)
	// TODO

	// Set up request
	ep := this.endpoint.URL
	ep.Path = "/write"
	req, err := http.NewRequest(http.MethodPost, ep.String(), buffer)
	if err != nil {
		return fmt.Errorf("%w: %q", err, ep.String())
	}
	req.Header.Set("Content-Type", "text/plain")

	// Set credentials
	if this.endpoint.user != "" {
		req.SetBasicAuth(this.endpoint.user, this.endpoint.password)
	}

	// Set parameters
	params := req.URL.Query()
	if db := this.endpoint.db; db != "" {
		params.Set("db", db)
	}
	req.URL.RawQuery = params.Encode()

	// Perform the request
	resp, err := this.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check status of request
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	} else if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v: %q", strings.TrimSpace(string(body)), ep.String())
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Writer) String() string {
	str := "<influxdb writer"
	str += " endpoint=" + fmt.Sprint(this.endpoint)
	if this.version != "" {
		str += " version=" + strconv.Quote(this.version)
	}
	return str + ">"
}

func (e endpoint) String() string {
	str := "<endpoint"
	str += " url=" + strconv.Quote(e.URL.String())
	if e.db != "" {
		str += " db=" + strconv.Quote(e.db)
	}
	if e.user != "" {
		str += " user=" + strconv.Quote(e.user)
	}
	if e.password != "" {
		str += " password=" + strings.Repeat("*", len(e.password))
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Parse -influxdb.url parameter to extract host,port,username,password,database
func parseUrl(value string) (endpoint, error) {
	// Set default value
	if value == "" {
		value = DefaultEndpoint
	}
	// Parse URL
	u, err := url.Parse(value)
	if err != nil {
		return endpoint{}, err
	}
	// Check empty host when no scheme set
	if u.Scheme == "" && u.Host == "" {
		arr := strings.SplitN(u.Path, "/", 2)
		u.Host = arr[0]
		if len(arr) > 1 {
			u.Path = "/" + arr[1]
		} else {
			u.Path = "/"
		}
	}
	// Check scheme
	if u.Scheme == "" {
		u.Scheme = DefaultScheme
	}
	// Check port
	if u.Port() == "" {
		u.Host = fmt.Sprintf("%s:%d", u.Host, DefaultPort)
	}
	// Make sure scheme is http or https
	if u.Scheme != "http" && u.Scheme != "https" {
		return endpoint{}, gopi.ErrBadParameter.WithPrefix("Unsupported scheme ", strconv.Quote(u.Scheme))
	} else if db, err := parseDatabase(u); err != nil {
		return endpoint{}, err
	} else {
		user, password := "", ""
		if u.User != nil {
			user = u.User.Username()
			password, _ = u.User.Password()
		}
		u.Path = "/"
		return endpoint{*u, db, user, password}, nil
	}
}

func parseDatabase(value *url.URL) (string, error) {
	db := strings.Trim(value.Path, "/")
	return db, nil
}
