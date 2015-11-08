package youtubeapi

import (
	"errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	//    "log"
)

type YouTubeService struct {
    service      *youtube.Service
	contentowner string
}

var (
	ErrorInvalidServiceAccount = errors.New("Invalid service account")
    ErrorInvalidClientSecrets  = errors.New("Invalid client secrets configuration")
	ErrorMissingContentOwner   = errors.New("Missing content owner")
	ErrorResponse              = errors.New("Bad Response")
)

// Returns a service object given service account details
func NewYouTubeServiceFromServiceAccountJSON(filename string, contentowner string) (*YouTubeService, error) {
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
	service, err := youtube.New(sa_config.Client(oauth2.NoContext))
	if err != nil {
		return nil, ErrorInvalidServiceAccount
	}
	return &YouTubeService{service, contentowner}, nil
}

// Returns a service object given client secrets details
func NewYouTubeServiceFromClientSecretsJSON(filename string) (*YouTubeService, error) {
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, ErrorInvalidClientSecrets
    }
    config, err := google.DefaultClient(oauth2.NoContext,youtube.YoutubeScope)
    if err != nil {
        return nil, ErrorInvalidClientSecrets
    }
    return &YouTubeService{ service,"" },nil
}

// Returns set of channel items for YouTube service. Can return several, in the
// case of service accounts, or a single one, based on simple OAuth authentication
func (ctx *YouTubeService) ChannelsList() (*youtube.ChannelListResponse, error) {
	call := ctx.service.Channels.List("contentDetails,snippet").OnBehalfOfContentOwner(ctx.contentowner).MaxResults(50).ManagedByMe(true)
	response, err := call.Do()
	if err != nil {
		return nil, ErrorResponse
	}
	return response,nil
}

