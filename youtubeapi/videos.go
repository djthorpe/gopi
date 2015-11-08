package youtubeapi

import (
	"google.golang.org/api/youtube/v3"
)

// Returns a list of videos for a playlist
func (ctx *YouTubeService) VideosList(part string) ([]*youtube.Video, error) {
	var call *youtube.VideosListCall
	if ctx.partnerapi {
		call = ctx.service.Videos.List(part).OnBehalfOfContentOwner(ctx.contentowner)
	} else {
		call = ctx.service.Videos.List(part)
	}
	response, err := call.MaxResults(50).Do()
	if err != nil {
		return nil, ErrorResponse
	}
	return response.Items,nil
}

