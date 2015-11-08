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
	response, _ := call.Do()
	if err != nil {
		return nil, ErrorResponse
	}
	return response.Items,nil
}

