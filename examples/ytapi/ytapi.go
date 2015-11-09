/* Copyright David Thorpe 2015 All Rights Reserved
   This package demonstrates calling the YouTube API
*/
package main

import (
	"flag"
	"os"
	"log"
	"os/user"
	"path/filepath"
	"strconv"
	"fmt"

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
	debug                  = flag.Bool("debug",false,"Debug flag")
)

var (
    operations = map[string]func(*youtubeapi.YouTubeService) {
        "videos": ListVideos,
        "channels": ListChannels,
    }
)

const (
	credentialsPathMode = 0700
	clientid = "973959355861.apps.googleusercontent.com"
)

////////////////////////////////////////////////////////////////////////////////

func userDir() (userDir string) {
	currentUser, _ := user.Current()
	userDir = currentUser.HomeDir
	return
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Set Usage function
	flag.Usage = func() {
        fmt.Fprintf(os.Stderr,"Usage of %s:\n",filepath.Base(os.Args[0]))
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
		service, err = youtubeapi.NewYouTubeServiceFromServiceAccountJSON(filepath.Join(credentialsPath,*serviceAccountFilename), *contentOwner,*debug)
    } else {
		service, err = youtubeapi.NewYouTubeServiceFromClientSecretsJSON(filepath.Join(credentialsPath,*clientsecretFilename),filepath.Join(credentialsPath,*tokenFilename),*debug)
    }
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Perform operation
    operations[opname](service)
}

func ListVideos(service *youtubeapi.YouTubeService) {
	// setup table
	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{ "channel_title", "video_id", "video_title", "video_status" })
	table.SetAutoFormatHeaders(false)

	// obtain channels
    channels, err := service.SetMaxResults(0).ChannelsList("contentDetails")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// obtain playlist items
	for _,channel := range channels {
		playlist := youtubeapi.YouTubePlaylistID(channel.ContentDetails.RelatedPlaylists.Uploads)
		videos, err := service.SetMaxResults(0).VideosForPlaylist("id,snippet,status",playlist)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		for _,video := range videos {
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


func ListChannels(service *youtubeapi.YouTubeService) {
    channels, err := service.ChannelsList("snippet,statistics")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{ "channel_id", "channel_title", "subscriber_count", "video_count","view_count" })
	table.SetAutoFormatHeaders(false)

	// Iterate through the channels
	for _, channel := range channels {
		table.Append([]string{
            channel.Id,
            channel.Snippet.Title,
			strconv.FormatUint(channel.Statistics.SubscriberCount,10),
			strconv.FormatUint(channel.Statistics.VideoCount,10),
			strconv.FormatUint(channel.Statistics.ViewCount,10),
        })
	}

	// Output the table
	table.Render()
}
