package youtubeapi

import (
	"google.golang.org/api/youtube/v3"
)

// Returns set of channel items for YouTube service. Can return several, in the
// case of service accounts, or a single one, based on simple OAuth authentication
func (this *YouTubeService) ChannelsList(part string) ([]*youtube.Channel, error) {
	var call *youtube.ChannelsListCall
	if this.partnerapi {
		call = this.service.Channels.List(part).OnBehalfOfContentOwner(this.contentowner).ManagedByMe(true)
	} else {
		call = this.service.Channels.List(part).Mine(true)
	}

    nextPageToken := ""
    items := make([]*youtube.Channel,0,this.maxresults)
    for {
        response, err := call.MaxResults(0).PageToken(nextPageToken).Do()
        if err != nil {
            return nil, ErrorResponse
        }
        items = append(items,response.Items...)
        nextPageToken = response.NextPageToken
        if nextPageToken == "" {
            break
        }
    }
	return items,nil
}

