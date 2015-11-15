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
			return nil, ErrorMissingChannelFlag
		}
	}

	// set status
	if this.status != "" {
		call = call.BroadcastStatus(this.status)
	} else if this.video != "" {
		call = call.Id(this.video)
	} else {
		call = call.BroadcastStatus("all")
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

func (this *YouTubeService) BindBroadcast(part string) ([]*youtube.LiveBroadcast, error) {
	// check for stream and video arguments
	if this.stream == "" || this.video == "" {
		return nil, ErrorMissingBindFlags
	}

	return this.BindUnbindBroadcast(part)
}

func (this *YouTubeService) UnbindBroadcast(part string) ([]*youtube.LiveBroadcast, error) {
	// check video argument
	if this.video == "" {
		return nil, ErrorMissingVideoFlag
	}

	return this.BindUnbindBroadcast(part)
}

func (this *YouTubeService) BindUnbindBroadcast(part string) ([]*youtube.LiveBroadcast, error) {

	// create call for bind
	call := this.service.LiveBroadcasts.Bind(this.video, part)

	// set authentication
	if this.partnerapi {
		call = call.OnBehalfOfContentOwner(this.contentowner)
		// check channel argument
		if this.channel != "" {
			call = call.OnBehalfOfContentOwnerChannel(this.channel)
		}
	}

	// set stream
	if this.stream != "" {
		call = call.StreamId(this.stream)
	}

	// make call
	response, err := call.Do()
	if err != nil {
		return nil, ErrorResponse
	}

	// return items
	return []*youtube.LiveBroadcast{response}, nil
}
