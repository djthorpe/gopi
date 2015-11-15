package youtubeapi

import (
	"google.golang.org/api/youtube/v3"
	"strings"
)

// Returns PlayListItems within a single playlist
func (this *YouTubeService) PlaylistItemsForPlaylist(part string, playlist YouTubePlaylistID) ([]*youtube.PlaylistItem, error) {
	var call *youtube.PlaylistItemsListCall
	if this.partnerapi {
		call = this.service.PlaylistItems.List(part).OnBehalfOfContentOwner(this.contentowner)
	} else {
		call = this.service.PlaylistItems.List(part)
	}

	// set the playlist flag
	call = call.PlaylistId(string(playlist))

	// initialize the paging
	var maxresults = this.maxresults
	nextPageToken := ""
	items := make([]*youtube.PlaylistItem, 0, this.maxresults)
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

// Returns Videos within a single playlist
func (this *YouTubeService) VideosForPlaylist(part string, playlist YouTubePlaylistID) ([]*youtube.Video, error) {
	var call *youtube.PlaylistItemsListCall
	if this.partnerapi {
		call = this.service.PlaylistItems.List("contentDetails").OnBehalfOfContentOwner(this.contentowner)
	} else {
		call = this.service.PlaylistItems.List("contentDetails")
	}

	// set the playlist flag
	call = call.PlaylistId(string(playlist))

	// Create the video call
	var videoCall *youtube.VideosListCall
	if this.partnerapi {
		videoCall = this.service.Videos.List(part).OnBehalfOfContentOwner(this.contentowner)
	} else {
		videoCall = this.service.Videos.List(part)
	}

	// page through the items
	var maxresults = this.maxresults
	var nextPageToken string = ""
	var items = make([]*youtube.Video, 0, maxresults)

	for {
		var pagingresults = int64(maxresults)
		if pagingresults <= 0 {
			pagingresults = YouTubeMaxPagingResults
		} else if pagingresults > YouTubeMaxPagingResults {
			pagingresults = YouTubeMaxPagingResults
		}
		response, err := call.MaxResults(pagingresults).PageToken(nextPageToken).Do()
		if err != nil {
			return nil, ErrorResponse
		}

		// Create array of VideoID items
		var playlistItems = make([]string, 0, len(response.Items))
		for _, item := range response.Items {
			playlistItems = append(playlistItems, item.ContentDetails.VideoId)
		}

		// return the videos
		videoResponse, err := videoCall.Id(strings.Join(playlistItems, ",")).MaxResults(int64(len(playlistItems))).Do()
		if err != nil {
			return nil, ErrorResponse
		}

		// append items
		items = append(items, videoResponse.Items...)

		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return items, nil
}
