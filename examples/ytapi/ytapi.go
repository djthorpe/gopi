/* Copyright David Thorpe 2015 All Rights Reserved
   This package demonstrates calling the YouTube API
*/
package main

import (
	"flag"
	"fmt"
	"log"
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
	statusFlag             = flag.String("status", "", "Status filter")
)

var (
	operations = map[string]func(*youtubeapi.YouTubeService){
		"videos":     ListVideos,     // --channel=<id> --maxresults=<n>
		"channels":   ListChannels,   // --channel=<id> --maxresults=<n>
		"broadcasts": ListBroadcasts, // --channel=<id> --maxresults=<n> --status=<active|all|completed|upcoming>
		"streams":    ListStreams,    // --channel=<id> --maxresults=<n>
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
	// Set status
	if err := service.SetStatus(*statusFlag); err != nil {
		log.Fatalf("Error with --status flag: %v\n", err)
	}

}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Set Usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", filepath.Base(os.Args[0]))
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

func ListVideos(service *youtubeapi.YouTubeService) {
	// setup table
	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"channeltitle", "video", "title", "privacy"})
	table.SetAutoFormatHeaders(false)

	// obtain channels
	channels, err := service.SetMaxResults(0).ChannelsList("contentDetails")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

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

func ListChannels(service *youtubeapi.YouTubeService) {
	channels, err := service.ChannelsList("snippet,statistics")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"channel", "title", "subscriber_count", "video_count", "view_count"})
	table.SetAutoFormatHeaders(false)

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

func ListBroadcasts(service *youtubeapi.YouTubeService) {
	broadcasts, err := service.BroadcastsList("snippet,status")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{ "broadcast", "title", "privacy", "status", "default" })
	table.SetAutoFormatHeaders(false)

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

func ListStreams(service *youtubeapi.YouTubeService) {
	streams, err := service.StreamsList("snippet,cdn,status")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{ "stream", "title", "format", "status", "default" })
	table.SetAutoFormatHeaders(false)

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

