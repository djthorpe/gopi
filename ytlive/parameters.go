package ytlive

import (
	"flag"
	"os"
	"encoding/csv"
	"log"
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////

var (
	inputFilename = flag.String("input","","Input file")
	flagChannel = flag.String("channel","","Channel name or identifier")
	flagVideo = flag.String("video","","Video identifier")
	flagStream = flag.String("stream","","Stream name or identifier")
//	flagContentOwner = flag.String("contentowner","","Content Owner identifier")
)

////////////////////////////////////////////////////////////////////////////////

type InputRecord struct {
	Channel string `csv:"channel,user,username"`
	Video string `csv:"video,video_id,videoid"`
	Stream string `csv:"stream,stream_id,streamname,stream_name"`
	Title string `csv:"title"`
	Description string `csv:"description"`
	Category string `csv:"category,category_id"`
	Keywords string `csv:"keywords"`
	StartTime string `csv:"start time,start"`
	Duration string `csv:"duration,end time,end"`
	Thumbnail string `csv:"thumbnail"`
}

////////////////////////////////////////////////////////////////////////////////
// Import parameters

func importFlatFile(filename *string) ([][]string, error) {
	fileHandle,err := os.Open(*filename)
	if err != nil {
		return nil,err
	}
	defer fileHandle.Close()
	reader := csv.NewReader(fileHandle)
	data,err := reader.ReadAll()
	if err != nil {
		return nil,err
	}
	return data,nil
}

func importFile(filename *string) (map[string]string, error) {
	// import data including the headers, etc
	data,err := importFlatFile(filename)
	if err != nil {
		return nil,err
	}
	if len(data) < 2 {
		return nil,fmt.Errorf("Empty file or missing data")
	}
	return nil,nil
}

////////////////////////////////////////////////////////////////////////////////
// Import parameters

func Import() ([]InputRecord, error) {
	// Set the defaults for the InputRecord
	defaults := InputRecord{
		Channel: *flagChannel,
		Video: *flagVideo,
		Stream: *flagStream,
	}

	log.Printf("Defaults = %v",defaults)

	// If there is an 'input' flag argument, then read the array
	if(inputFilename != nil) {
		data,err := importFile(inputFilename)
		if(err != nil) {
			return nil,err
		}
		log.Printf("Data = %v",data)
	} else {

	}

	// return success
	input := []InputRecord{
		InputRecord{ },
	}
	return input,nil
}
