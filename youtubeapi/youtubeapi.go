package youtubeapi

import (
	"os"
	"errors"
	"time"
	"fmt"
	"net/http"
	"net/http/httptest"
	"encoding/gob"
	"io/ioutil"
	"os/exec"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)


import "log"

// YouTubeService object which contains the main context for calling the
// YouTube API
type YouTubeService struct {
    service      *youtube.Service
	token        *oauth2.Token
	contentowner string
	partnerapi   bool
	debug        bool
}

var (
	ErrorInvalidServiceAccount = errors.New("Invalid service account")
    ErrorInvalidClientSecrets  = errors.New("Invalid client secrets configuration")
	ErrorMissingContentOwner   = errors.New("Missing content owner")
	ErrorCacheTokenRead        = errors.New("Invalid Cache Token")
	ErrorCacheTokenWrite       = errors.New("Unable to create cache token")
	ErrorTokenExchange         = errors.New("Token Exchange Error")
	ErrorResponse              = errors.New("Bad Response")
)

// Returns a service object given service account details
func NewYouTubeServiceFromServiceAccountJSON(filename string, contentowner string,debug bool) (*YouTubeService, error) {
	if len(contentowner) == 0 {
		return nil, ErrorMissingContentOwner
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, ErrorInvalidServiceAccount
	}
	sa_config, err := google.JWTConfigFromJSON(bytes, youtube.YoutubeScope, youtube.YoutubepartnerScope)
	if err != nil {
		return nil, ErrorInvalidServiceAccount
	}
	ctx := context.Background()
	service, err := youtube.New(sa_config.Client(ctx))
	if err != nil {
		return nil, ErrorInvalidServiceAccount
	}
	return &YouTubeService{ service, nil, contentowner, true, debug }, nil
}

// Returns a service object given client secrets details
func NewYouTubeServiceFromClientSecretsJSON(clientsecrets string,tokencache string,debug bool) (*YouTubeService, error) {
    bytes, err := ioutil.ReadFile(clientsecrets)
    if err != nil {
        return nil, ErrorInvalidClientSecrets
    }
	config,err := google.ConfigFromJSON(bytes,youtube.YoutubeScope)
	if err != nil {
        return nil, ErrorInvalidClientSecrets
	}
	ctx := getContext(debug)

	// Attempt to get token from cache
	token, err := tokenFromFile(tokencache)
	if err != nil {
		token,err = tokenFromWeb(config,ctx)
		saveToken(tokencache,token)
	}
	if err != nil {
		return nil, ErrorInvalidClientSecrets
	}

	// create client
	service,err := youtube.New(config.Client(ctx,token))
	if err != nil {
		return nil, ErrorInvalidClientSecrets
	}

    return &YouTubeService{ service, token, "", false, debug },nil
}

// Returns context
func getContext(debug bool) (context.Context) {
	ctx := context.Background()
	if(debug) {
		ctx = context.WithValue(ctx,oauth2.HTTPClient,&http.Client{
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
func saveToken(filename string,token *oauth2.Token) (error) {
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("Error writing to %s",filename)
		return ErrorCacheTokenWrite
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(token)
	return nil
}

// Creates a webserver for user interaction with Google
func tokenFromWeb(config *oauth2.Config,ctx context.Context) (*oauth2.Token,error) {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			http.Error(rw,"State doesn't match",500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw,"<h1>Success</h1>Authorized - You can now close this window")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		http.Error(rw,"No Code",500)
	}))
	defer ts.Close()
	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)
	go openURL(authURL)
	code := <-ch
	token, err := config.Exchange(ctx,code)
	if err != nil {
		return nil,ErrorTokenExchange
	}
	return token,nil
}

func openURL(url string) {
	try := []string{"xdg-open", "google-chrome", "open"}
	for _, bin := range try {
		err := exec.Command(bin, url).Run()
		if err == nil {
			return
		}
	}
}

// Save token
func (ctx *YouTubeService) SaveCredentials(filename string) (error) {
	return nil
}
