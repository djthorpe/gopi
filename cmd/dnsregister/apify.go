package main

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *app) RegisterExternalAddress(addr net.IP, host, user, passwd string) (string, error) {
	// Create request
	req, err := http.NewRequest("GET", uriGoogle, nil)
	if err != nil {
		return "", err
	}

	// Set up the request
	values := req.URL.Query()
	values.Set("hostname", host)
	values.Set("ip", addr.String())
	req.URL.RawQuery = values.Encode()
	req.URL.User = url.UserPassword(user, passwd)

	// Get response
	response, err := this.Do(req)
	if err != nil {
		return "", err
	}
	// Check response parameters
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %v", gopi.ErrUnexpectedResponse, response.Status)
	} else if contentType, err := GetContentType(response); err != nil {
		return "", err
	} else if contentType != "text/plain" {
		return "", fmt.Errorf("%w: %q", gopi.ErrUnexpectedResponse, contentType)
	}
	// Read body
	defer response.Body.Close()
	if data, err := ioutil.ReadAll(response.Body); err != nil {
		return "", err
	} else if status := strings.SplitN(string(data), " ", 2); len(status) < 1 {
		return string(data), nil
	} else {
		return status[0], nil
	}
}

func (this *app) GetExternalAddress() (net.IP, error) {
	// Perform request
	req, err := http.NewRequest("GET", uriApify, nil)
	if err != nil {
		return nil, err
	}
	// Get response
	response, err := this.Do(req)
	if err != nil {
		return nil, err
	}
	// Check response parameters
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %v", gopi.ErrUnexpectedResponse, response.Status)
	} else if contentType, err := GetContentType(response); err != nil {
		return nil, err
	} else if contentType != "text/plain" {
		return nil, fmt.Errorf("%w: %q", gopi.ErrUnexpectedResponse, contentType)
	}
	// Read body
	defer response.Body.Close()
	if data, err := ioutil.ReadAll(response.Body); err != nil {
		return nil, err
	} else if ip := net.ParseIP(string(data)); ip == nil {
		return nil, fmt.Errorf("%w: %q", gopi.ErrUnexpectedResponse, string(data))
	} else {
		return ip, nil
	}
}

func GetContentType(resp *http.Response) (string, error) {
	if v := resp.Header.Get("Content-Type"); v == "" {
		return "application/octet-stream", nil
	} else if t, _, err := mime.ParseMediaType(v); err != nil {
		return "", err
	} else {
		return t, nil
	}
}
