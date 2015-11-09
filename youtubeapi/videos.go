package youtubeapi

import (
	"strings"
	"google.golang.org/api/youtube/v3"
)

// Returns PlayListItems within a single playlist
func (this *YouTubeService) PlaylistItemsForPlaylist(part string,playlist YouTubePlaylistID) ([]*youtube.PlaylistItem, error) {
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
    items := make([]*youtube.PlaylistItem,0,this.maxresults)
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
        items = append(items,response.Items...)
        nextPageToken = response.NextPageToken
        if nextPageToken == "" {
            break
        }
    }
	return items,nil
}

// Returns Videos within a single playlist
func (this *YouTubeService) VideosForPlaylist(part string,playlist YouTubePlaylistID) ([]*youtube.Video, error) {

	// Fetch Playlist Items
	playlistItems,err := this.PlaylistItemsForPlaylist("contentDetails",playlist)
	if err != nil {
		return nil,err
	}

	// Generate List of VideoID's
	videoIdList := make([]string,0,len(playlistItems))
	for _,video := range playlistItems {
		videoIdList = append(videoIdList,video.ContentDetails.VideoId)
	}

	// Generate the call to youtube.videos.list
	var call *youtube.VideosListCall
	if this.partnerapi {
		call = this.service.Videos.List(part).OnBehalfOfContentOwner(this.contentowner)
	} else {
		call = this.service.Videos.List(part)
	}
	call = call.Id(strings.Join(videoIdList,","))

	// Paging through calls
    nextPageToken := ""
    items := make([]*youtube.Video,0,len(videoIdList))
    for {
        response, err := call.MaxResults(YouTubeMaxPagingResults).PageToken(nextPageToken).Do()
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


