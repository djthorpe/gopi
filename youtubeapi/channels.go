package youtubeapi

import (
	"google.golang.org/api/youtube/v3"
)

////////////////////////////////////////////////////////////////////////////////

// Returns set of channel items for YouTube service. Can return several, in the
// case of service accounts, or a single one, based on simple OAuth authentication
func (this *YouTubeService) ChannelsList(part string) ([]*youtube.Channel, error) {

	// create call for channels
	call := this.service.Channels.List(part)

	// set channel or channels
	if this.channel != "" {
		call = call.Id(this.channel)
	} else if this.partnerapi {
		call = call.OnBehalfOfContentOwner(this.contentowner).ManagedByMe(true)
	} else {
		call = call.Mine(true)
	}

	// page through results
	var maxresults = this.maxresults
	nextPageToken := ""
	items := make([]*youtube.Channel, 0, maxresults)
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
