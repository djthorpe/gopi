package youtubeapi

import (
	"google.golang.org/api/youtube/v3"
)

////////////////////////////////////////////////////////////////////////////////

// Returns set of broadcast items for YouTube service
func (this *YouTubeService) BroadcastsList(part string) ([]*youtube.LiveBroadcast, error) {

	// create call for broadcasts
	call := this.service.LiveBroadcasts.List(part)

	// set channel
	if this.partnerapi {
		call = call.OnBehalfOfContentOwner(this.contentowner)
		// check channel argument
		if this.channel != "" {
			call = call.OnBehalfOfContentOwnerChannel(this.channel)
		} else {
			return nil,ErrorMissingChannelFlag
		}
	}

	// set status
	if this.status != "" {
		call = call.BroadcastStatus(this.status)
	}

	// page through results
	var maxresults = this.maxresults
	nextPageToken := ""
	items := make([]*youtube.LiveBroadcast, 0, maxresults)
	for {
		var pagingresults = int64(maxresults) - int64(len(items))
		if pagingresults <= 0 {
			pagingresults = YouTubeMaxPagingResults
		} else if pagingresults > YouTubeMaxPagingResults {
			pagingresults = YouTubeMaxPagingResults
		}
		response, err := call.MaxResults(pagingresults).PageToken(nextPageToken).Do()
		if err != nil {
			return nil, ErrorResponse
		}
		items = append(items, response.Items...)
		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return items, nil
}

////////////////////////////////////////////////////////////////////////////////

