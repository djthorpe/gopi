package youtubeapi

import (
	"google.golang.org/api/youtube/v3"
)

// Returns set of channel items for YouTube service. Can return several, in the
// case of service accounts, or a single one, based on simple OAuth authentication
func (ctx *YouTubeService) ChannelsList(part string) ([]*youtube.Channel, error) {
	var call *youtube.ChannelsListCall
	if ctx.partnerapi {
		call = ctx.service.Channels.List(part).OnBehalfOfContentOwner(ctx.contentowner).ManagedByMe(true)
	} else {
		call = ctx.service.Channels.List(part).Mine(true)
	}
	response, err := call.MaxResults(50).Do()
	if err != nil {
		return nil, ErrorResponse
	}
	return response.Items,nil
}

