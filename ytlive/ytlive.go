package main

import (
	"flag"
	"os"
//	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

var (
	clientsecretFilename = flag.String("clientsecret","client_secret.json","Client secret filename")
	serviceAccountFilename = flag.String("serviceaccount","service_account.json","Service account filename")
	tokenFilename = flag.String("authtoken","oauth_token","OAuth token filename")
	credentialsFolder = flag.String("credentials",".credentials","Folder containing credentials")
	contentOwner = flag.String("contentowner","","Content Owner ID")
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
	credentialsPath := filepath.Join(userDir(),*credentialsFolder)
	if credentialsPathInfo, err := os.Stat(credentialsPath); err != nil || !credentialsPathInfo.IsDir() {
		log.Fatalf("Missing path: %v\n",credentialsPathInfo)
	}

	// Read service account
	bytes, err := ioutil.ReadFile(filepath.Join(credentialsPath,*serviceAccountFilename))
	if err != nil {
		log.Fatalf("Invalid service account file: %s\n",*serviceAccountFilename)
	}
	sa_config,err := google.JWTConfigFromJSON(bytes,youtube.YoutubeScope,youtube.YoutubepartnerScope)
	if err != nil {
		log.Fatalf("Invalid service account configuration: %s\n",serviceAccountFilename)
	}

	// create a youtube partner service
	client := sa_config.Client(oauth2.NoContext)
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// List channels
	call := service.Channels.List("contentDetails,snippet").OnBehalfOfContentOwner(*contentOwner).MaxResults(50).ManagedByMe(true)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	// Create table writer object
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Channel","Uploads"})
	table.SetAutoFormatHeaders(false)

	// Iterate through the channels
	for _, channel := range response.Items {
		table.Append([]string{ channel.Snippet.Title, channel.ContentDetails.RelatedPlaylists.Uploads })
	}

	// Output the table
	table.Render()

	/*

		playlistId := channel.ContentDetails.RelatedPlaylists.Uploads
		// Print the playlist ID for the list of uploaded videos.
		fmt.Printf("Videos in list %s\r\n", playlistId)
		nextPageToken := ""
		for {
			// Call the playlistItems.list method to retrieve the
			// list of uploaded videos. Each request retrieves 50
			// videos until all videos have been retrieved.
			playlistCall := service.PlaylistItems.List("snippet").PlaylistId(playlistId).MaxResults(50).PageToken(nextPageToken)

			playlistResponse, err := playlistCall.Do()

			if err != nil {
				// The playlistItems.list method call returned an error.
				log.Fatalf("Error fetching playlist items: %v", err.Error())
			}

			for _, playlistItem := range playlistResponse.Items {
				title := playlistItem.Snippet.Title
				videoId := playlistItem.Snippet.ResourceId.VideoId
				fmt.Printf("  %v, (%v)\r\n", title, videoId)
			}

			// Set the token to retrieve the next page of results
			// or exit the loop if all results have been retrieved.
			nextPageToken = playlistResponse.NextPageToken
			if nextPageToken == "" {
				break
			}
			fmt.Println()
		}
	}
	*/
}

/*
	// Read client secrets
	bytes, err = ioutil.ReadFile(filepath.Join(credentialsPath,*clientsecretFilename))
	if err != nil {
		log.Fatalf("Invalid client secrets file: %s\n",clientsecretFilename)
	}

	_, err := google.ConfigFromJSON(bytes,youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Invalid client secrets configuration: %v\n",err)
	}
*/
