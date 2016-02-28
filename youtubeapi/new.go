package youtubeapi

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
)

// Service object which contains the main context for calling the YouTube API
type Service struct {
	service      *youtube.Service
	token        *oauth2.Token
	contentowner string
	partnerapi   bool
	debug        bool
}

// Parameters object contains parameters used when calling the API
type Parameters struct {
    channel      string
    video        string
    broadcast    string
    stream       string
    status       string
    maxresults   uint
}

// YouTube Identifiers
type YouTubePlaylistID string
type YouTubeVideoID string
type YouTubeChannelID string
type YouTubeStreamID string

// Constants
const (
	YouTubeMaxPagingResults = 50
)

// Returns a service object given service account details
func NewYouTubeServiceFromServiceAccountJSON(filename string, contentowner string, debug bool) (*Service, error) {
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
	ctx := getContext(debug)
	service, err := youtube.New(sa_config.Client(ctx))
	if err != nil {
		return nil, ErrorInvalidServiceAccount
	}
	this := new(Service)
	this.service = service
	this.contentowner = contentowner
	this.partnerapi = true
	this.debug = debug
	return this, nil
}

// Returns a service object given client secrets details
func NewYouTubeServiceFromClientSecretsJSON(clientsecrets string, tokencache string, debug bool) (*Service, error) {
	bytes, err := ioutil.ReadFile(clientsecrets)
	if err != nil {
		return nil, ErrorInvalidClientSecrets
	}
	config, err := google.ConfigFromJSON(bytes, youtube.YoutubeScope)
	if err != nil {
		return nil, ErrorInvalidClientSecrets
	}
	ctx := getContext(debug)

	// Attempt to get token from cache
	token, err := tokenFromFile(tokencache)
	if err != nil {
		token, err = tokenFromWeb(config, ctx)
		saveToken(tokencache, token)
	}
	if err != nil {
		return nil, ErrorInvalidClientSecrets
	}

	// create client
	service, err := youtube.New(config.Client(ctx, token))
	if err != nil {
		return nil, ErrorInvalidClientSecrets
	}

	this := new(Service)
	this.service = service
	this.token = token
	this.partnerapi = false
	this.debug = debug
	return this, nil
}
