
package youtubeapi

import (
	"fmt"
	"time"
	"os"
	"os/exec"
	"encoding/gob"
	"net/http"
	"net/http/httptest"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Returns context
func getContext(debug bool) context.Context {
	ctx := context.Background()
	if debug {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
			Transport: &logTransport{http.DefaultTransport},
		})
	}
	return ctx
}

// Returns token from cache
func tokenFromFile(filename string) (*oauth2.Token, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, ErrorCacheTokenRead
	}
	t := new(oauth2.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

// Saves token
func saveToken(filename string, token *oauth2.Token) error {
	f, err := os.Create(filename)
	if err != nil {
		return ErrorCacheTokenWrite
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(token)
	return nil
}

// Creates a webserver for user interaction with Google
func tokenFromWeb(config *oauth2.Config, ctx context.Context) (*oauth2.Token, error) {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			http.Error(rw, "State doesn't match", 500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized - You can now close this window")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		http.Error(rw, "No Code", 500)
	}))
	defer ts.Close()
	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)
	go openURL(authURL)
	code := <-ch
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, ErrorTokenExchange
	}
	return token, nil
}

// Attempt to open a URL using a browser
func openURL(url string) {
	try := []string{"xdg-open", "google-chrome", "open"}
	for _, bin := range try {
		err := exec.Command(bin, url).Run()
		if err == nil {
			return
		}
	}
}
