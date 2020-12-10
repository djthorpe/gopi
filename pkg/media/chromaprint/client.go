package chromaprint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Client struct {
	gopi.Unit
	gopi.Logger
	sync.Mutex
	*http.Client
	*url.URL

	last time.Time
}

type Meta uint

type results struct {
	Status  string   `json:"status"`
	Error   reason   `json:"error"`
	Results []*Match `json:"results"`
}

type reason struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Register a clientId: https://acoustid.org/login
	defaultClientId = "xAJ7ALREvAU"

	// Timeout requests after 15 seconds
	defaultTimeout = 15 * time.Second

	// The API endpoint
	baseUrl = "https://api.acoustid.org/v2"

	// defaultQps rate limits number of requests per second
	defaultQps = 3
)

const (
	META_NONE      Meta = 0
	META_RECORDING Meta = (1 << iota)
	META_RECORDINGID
	META_RELEASE
	META_RELEASEID
	META_RELEASEGROUP
	META_RELEASEGROUPID
	META_TRACK
	META_COMPRESS
	META_USERMETA
	META_SOURCE
	META_MIN = META_RECORDING
	META_MAX = META_SOURCE
	META_ALL = META_RECORDING | META_RECORDINGID | META_RELEASE | META_RELEASEID | META_RELEASEGROUP | META_RELEASEGROUPID | META_TRACK | META_COMPRESS | META_USERMETA | META_SOURCE
)

var (
	ErrQueryRateExceeded = fmt.Errorf("Query Rate Exceeded (%v qps)", defaultQps)
)

////////////////////////////////////////////////////////////////////////////////
// NEW

func (this *Client) Define(cfg gopi.Config) {
	cfg.FlagDuration("acoustid.timeout", 0, "AcoustId client timeout")
	cfg.FlagString("acoustid.clientid", "", "AcoustId client identifier")
}

func (this *Client) New(cfg gopi.Config) error {
	timeout := cfg.GetDuration("acoustid.timeout")
	clientId := cfg.GetString("acoustid.clientid")

	// Set parameters
	if timeout == 0 {
		timeout = defaultTimeout
	}
	if clientId == "" {
		clientId = defaultClientId
	}
	this.Client = &http.Client{
		Timeout: timeout,
	}
	if url, err := url.Parse(baseUrl); err != nil {
		return err
	} else {
		v := url.Query()
		v.Set("client", clientId)
		url.RawQuery = v.Encode()
		this.URL = url
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Client) String() string {
	str := "<chromaprint.client"
	if u := this.URL; u != nil {
		str += " url=" + u.String()
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// LOOKUP

func (this *Client) Lookup(fingerprint string, duration time.Duration, flags Meta) ([]*Match, error) {
	// Check incoming parameters
	if fingerprint == "" || duration == 0 || flags == META_NONE {
		return nil, gopi.ErrBadParameter.WithPrefix("Lookup")
	}

	// Check Qps
	if this.last.IsZero() {
		if time.Since(this.last) < (time.Second / defaultQps) {
			return nil, ErrQueryRateExceeded
		}
	}

	// Set URL parameters
	params := url.Values{}
	params.Set("fingerprint", fingerprint)
	params.Set("duration", fmt.Sprint(duration.Seconds()))
	params.Set("meta", flags.String())

	url := this.requestUrl("lookup", params)
	if url == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("Lookup")
	}

	// Perform request
	now := time.Now()
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	response, err := this.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Debug output
	if this.Logger != nil {
		this.Debug("Lookup", url, "took", time.Since(now), "returned", response.Status)
	}

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Decode response body
	var r results
	if mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-type")); err != nil {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(err)
	} else if mimeType != "application/json" {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(mimeType)
	} else if err := json.Unmarshal(body, &r); err != nil {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix(err)
	}

	// Check for errors
	if r.Status != "ok" {
		return nil, fmt.Errorf("%w: %v (code %v)", gopi.ErrBadParameter, r.Error.Message, r.Error.Code)
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %v (code %v)", gopi.ErrBadParameter, response.Status, response.StatusCode)
	}

	// Set response time for calculating qps
	this.last = now

	// Return success
	return r.Results, nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Client) requestUrl(path string, v url.Values) *url.URL {
	url, err := url.Parse(this.URL.String())
	if err != nil {
		return nil
	}
	// Copy params
	params := this.URL.Query()
	for k := range v {
		params[k] = v[k]
	}
	url.RawQuery = params.Encode()

	// Set path
	url.Path = filepath.Join(url.Path, path)

	// Return URL
	return url
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (m Meta) String() string {
	if m == META_NONE {
		return m.FlagString()
	}
	str := ""
	for v := META_MIN; v <= META_MAX; v <<= 1 {
		if m&v == v {
			str += v.FlagString() + ","
		}
	}
	return strings.TrimSuffix(str, ",")
}

func (m Meta) FlagString() string {
	switch m {
	case META_NONE:
		return ""
	case META_RECORDING:
		return "recordings"
	case META_RECORDINGID:
		return "recordingids"
	case META_RELEASE:
		return "releases"
	case META_RELEASEID:
		return "releaseids"
	case META_RELEASEGROUP:
		return "releasegroups"
	case META_RELEASEGROUPID:
		return "releasegroupids"
	case META_TRACK:
		return "tracks"
	case META_COMPRESS:
		return "compress"
	case META_USERMETA:
		return "usermeta"
	case META_SOURCE:
		return "sources"
	default:
		return "[?? Invalid Meta value]"
	}
}
