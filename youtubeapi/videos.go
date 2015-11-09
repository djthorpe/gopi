package youtubeapi

import (
	"google.golang.org/api/youtube/v3"
)

// Returns set of video items for YouTube service, given one or more video ids.
func (this *YouTubeService) VideosForPlaylist(part string,playlist YouTubePlaylistID) ([]*youtube.PlaylistItem, error) {
	var call *youtube.PlaylistItemsListCall
	if this.partnerapi {
		call = this.service.PlaylistItems.List(part).OnBehalfOfContentOwner(this.contentowner)
	} else {
		call = this.service.PlaylistItems.List(part)
	}

	// set the playlist flag
	call = call.PlaylistId(string(playlist))

    nextPageToken := ""
    items := make([]*youtube.PlaylistItem,0,this.maxresults)
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
