package youtubeapi

import (
	"google.golang.org/api/youtube/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TODO Cache results of keys => streams
// Returns stream for key
func (this *YouTubeService) StreamForKey(key string) (string, error) {
	call := this.service.LiveStreams.List("id")

	// set channel
	if this.partnerapi {
		call = call.OnBehalfOfContentOwner(this.contentowner)
		// check channel argument
		if this.channel != "" {
			call = call.OnBehalfOfContentOwnerChannel(this.channel)
		} else {
			return "", ErrorMissingChannelFlag
		}
	} else {
		call = call.Mine(true)
	}

	_, err := call.Id(key).Do()
	if err != nil {
		return "", err
	}

	return key, nil
}

// Returns set of stream items for YouTube service
func (this *YouTubeService) StreamsList(part string) ([]*youtube.LiveStream, error) {

	// create call for broadcasts
	call := this.service.LiveStreams.List(part)

	// set channel
	if this.partnerapi {
		call = call.OnBehalfOfContentOwner(this.contentowner)
		// check channel argument
		if this.channel != "" {
			call = call.OnBehalfOfContentOwnerChannel(this.channel)
		} else {
			return nil, ErrorMissingChannelFlag
		}
	} else {
		call = call.Mine(true)
	}

	// page through results
	var maxresults = this.maxresults
	nextPageToken := ""
	items := make([]*youtube.LiveStream, 0, maxresults)
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
