package main

import (
	"flag"
	"os"
	//	"fmt"
	"log"
	"os/user"
	"path/filepath"

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
)

const (
	credentialsPathMode = 0700
)

////////////////////////////////////////////////////////////////////////////////

func userDir() (userDir string) {
	currentUser, _ := user.Current()
	userDir = currentUser.HomeDir
	return
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Read flags
	flag.Parse()

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
		service, err = youtubeapi.NewYouTubeServiceFromServiceAccountJSON(filepath.Join(credentialsPath,*serviceAccountFilename), *contentOwner)
    } else {
		service, err = youtubeapi.NewYouTubeServiceFromClientSecretsJSON(filepath.Join(credentialsPath,*clientsecretFilename))
    }
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

    response, err := service.ChannelsList()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Channel", "Uploads"})
	table.SetAutoFormatHeaders(false)
	// Iterate through the channels
	for _, channel := range response.Items {
		table.Append([]string{channel.Snippet.Title, channel.ContentDetails.RelatedPlaylists.Uploads})
	}
	// Output the table
	table.Render()
}
