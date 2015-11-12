/* Copyright David Thorpe 2015 All Rights Reserved
   This package demonstrates calling the YouTube API
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/djthorpe/gopi/youtubeapi"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

var (
	clientsecretFilename   = flag.String("clientsecret", "client_secret.json", "Client secret filename")
	serviceAccountFilename = flag.String("serviceaccount", "service_account.json", "Service account filename")
	tokenFilename          = flag.String("authtoken", "oauth_token", "OAuth token filename")
	credentialsFolder      = flag.String("credentials", ".credentials", "Folder containing credentials")
	contentOwner           = flag.String("contentowner", "", "Content Owner ID")
	debug                  = flag.Bool("debug", false, "Debug flag")
	channelFlag            = flag.String("channel", "", "Channel ID")
	videoFlag              = flag.String("video", "", "Video ID")
	streamFlag             = flag.String("stream", "", "Stream Key or ID")
	statusFlag             = flag.String("status", "", "Status filter")
)

var (
	operations = map[string]func(*youtubeapi.YouTubeService){
		"videos":     ListVideos,     // --channel=<id> --maxresults=<n>
		"channels":   ListChannels,   // --channel=<id> --maxresults=<n>
		"broadcasts": ListBroadcasts, // --channel=<id> --maxresults=<n> --status=<active|all|completed|upcoming>
		"streams":    ListStreams,    // --channel=<id> --maxresults=<n>
		"bind":       BindBroadcast,  // --video=<id> --stream=<key>
		"unbind":     UnbindBroadcast,// --video=<id>
	}
)

const (
	credentialsPathMode = 0700
	clientid            = "973959355861.apps.googleusercontent.com"
)

////////////////////////////////////////////////////////////////////////////////

func userDir() (userDir string) {
	currentUser, _ := user.Current()
	userDir = currentUser.HomeDir
	return
}

func setDefaults(service *youtubeapi.YouTubeService) {

	// Set channel
	if err := service.SetChannel(*channelFlag); err != nil {
		log.Fatalf("Error with --channel flag: %v\n", err)
	}

	// Set video
	if err := service.SetVideo(*videoFlag); err != nil {
		log.Fatalf("Error with --video flag: %v\n", err)
	}

	// Set stream
	if err := service.SetStream(*streamFlag); err != nil {
		log.Fatalf("Error with --stream flag: %v\n", err)
	}

	// Set status
	if err := service.SetStatus(*statusFlag); err != nil {
		log.Fatalf("Error with --status flag: %v\n", err)
	}

}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// enumerate operations
	operation_keys := make([]string,0,len(operations))
	for key := range operations {
	    operation_keys = append(operation_keys,key)
	}

	// Set usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\n\t%s <flags> <%v>\n\n",filepath.Base(os.Args[0]),strings.Join(operation_keys,"|"))
		fmt.Fprintf(os.Stderr, "Where <flags> are one or more of:\n\n")
		flag.PrintDefaults()
	}

	// Read flags, exit with no operation
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	opname := flag.Arg(0)
	if operations[opname] == nil {
		flag.Usage()
		os.Exit(1)
	}

	// Obtain path for credentials
	credentialsPath := filepath.Join(userDir(), *credentialsFolder)
	if credentialsPathInfo, err := os.Stat(credentialsPath); err != nil || !credentialsPathInfo.IsDir() {
		// if path is missing, try and create the folder
		if err := os.Mkdir(credentialsPath, credentialsPathMode); err != nil {
			log.Fatalf("Missing credentials folder: %v\n", credentialsPath)
		}
	}

	// If we have a content owner, then assume we're going to create the service
	// using a service account
	var service *youtubeapi.YouTubeService
	var err error
	if len(*contentOwner) > 0 {
		service, err = youtubeapi.NewYouTubeServiceFromServiceAccountJSON(filepath.Join(credentialsPath, *serviceAccountFilename), *contentOwner, *debug)
	} else {
		service, err = youtubeapi.NewYouTubeServiceFromClientSecretsJSON(filepath.Join(credentialsPath, *clientsecretFilename), filepath.Join(credentialsPath, *tokenFilename), *debug)
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Set defaults for flags
	setDefaults(service)

	// Perform operation
	operations[opname](service)
}

////////////////////////////////////////////////////////////////////////////////
// Videos

func ListVideos(service *youtubeapi.YouTubeService) {

	// obtain channels
	channels, err := service.SetMaxResults(0).ChannelsList("contentDetails")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create table writer object
	table := NewVideosTable

	// obtain playlist items
	for _, channel := range channels {
		playlist := youtubeapi.YouTubePlaylistID(channel.ContentDetails.RelatedPlaylists.Uploads)
		videos, err := service.SetMaxResults(0).VideosForPlaylist("id,snippet,status", playlist)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		for _, video := range videos {
			table.Append([]string{
				video.Snippet.ChannelTitle,
				video.Id,
				video.Snippet.Title,
				video.Status.PrivacyStatus,
			})
		}
	}

	// Output the table
	table.Render()
}

////////////////////////////////////////////////////////////////////////////////
// Channels

func ListChannels(service *youtubeapi.YouTubeService) {
	channels, err := service.ChannelsList("snippet,statistics")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create table writer object
	table := NewChannelsTable()

	// Iterate through the channels
	for _, channel := range channels {
		table.Append([]string{
			channel.Id,
			channel.Snippet.Title,
			strconv.FormatUint(channel.Statistics.SubscriberCount, 10),
			strconv.FormatUint(channel.Statistics.VideoCount, 10),
			strconv.FormatUint(channel.Statistics.ViewCount, 10),
		})
	}

	// Output the table
	table.Render()
}

////////////////////////////////////////////////////////////////////////////////
// Broadcasts

func ListBroadcasts(service *youtubeapi.YouTubeService) {
	broadcasts, err := service.BroadcastsList("snippet,status")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	table := NewBroadcastsTable()

	// Iterate through the broadcasts
	for _, broadcast := range broadcasts {
		table.Append([]string{
			broadcast.Id,
			broadcast.Snippet.Title,
			broadcast.Status.PrivacyStatus,
			broadcast.Status.LifeCycleStatus,
			strconv.FormatBool(broadcast.Snippet.IsDefaultBroadcast),
		})
	}

	// Output the table
	table.Render()
}

func BindBroadcast(service *youtubeapi.YouTubeService) {
	broadcasts, err := service.BindBroadcast("snippet,status")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	table := NewBroadcastsTable()

	// Iterate through the broadcasts
	for _, broadcast := range broadcasts {
		table.Append([]string{
			broadcast.Id,
			broadcast.Snippet.Title,
			broadcast.Status.PrivacyStatus,
			broadcast.Status.LifeCycleStatus,
			strconv.FormatBool(broadcast.Snippet.IsDefaultBroadcast),
		})
	}

	// Output the table
	table.Render()
}

func UnbindBroadcast(service *youtubeapi.YouTubeService) {
	broadcasts, err := service.UnbindBroadcast("snippet,status")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	table := NewBroadcastsTable()

	// Iterate through the broadcasts
	for _, broadcast := range broadcasts {
		table.Append([]string{
			broadcast.Id,
			broadcast.Snippet.Title,
			broadcast.Status.PrivacyStatus,
			broadcast.Status.LifeCycleStatus,
			strconv.FormatBool(broadcast.Snippet.IsDefaultBroadcast),
		})
	}

	// Output the table
	table.Render()
}

////////////////////////////////////////////////////////////////////////////////
// Streams

func ListStreams(service *youtubeapi.YouTubeService) {
	streams, err := service.StreamsList("snippet,cdn,status")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create table writer object
	table := NewStreamsTable()

	// Iterate through the broadcasts
	for _, stream := range streams {
		table.Append([]string{
			stream.Cdn.IngestionInfo.StreamName,
			stream.Snippet.Title,
			stream.Cdn.Format,
			stream.Status.StreamStatus,
			strconv.FormatBool(stream.Snippet.IsDefaultStream),
		})
	}

	// Output the table
	table.Render()
}

////////////////////////////////////////////////////////////////////////////////
// Create table objects for each type of return

func NewBroadcastsTable() (*tablewriter.Table) {
	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{ "broadcast", "title", "privacy", "status", "default" })
	table.SetAutoFormatHeaders(false)
	return table
}

func NewStreamsTable() (*tablewriter.Table) {
	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{ "stream", "title", "format", "status", "default" })
	table.SetAutoFormatHeaders(false)
	return table
}

func NewVideosTable() (*tablewriter.Table) {
	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"channeltitle", "video", "title", "privacy"})
	table.SetAutoFormatHeaders(false)
	return table
}

func NewChannelsTable() (*tablewriter.Table) {
	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"channel", "title", "subscriber_count", "video_count", "view_count"})
	table.SetAutoFormatHeaders(false)
	return table
}




