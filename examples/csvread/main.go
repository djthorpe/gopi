package main

import (
    "github.com/djthorpe/gopi/csvdata"
    "flag"
    "log"
    "os"
    "encoding/csv"
    "strings"
    "fmt"
)

////////////////////////////////////////////////////////////////////////////////

type PrivacyValue string

func (this *PrivacyValue) String() string {
    return strings.ToUpper(string(*this))
}

func (this *PrivacyValue) Set(sval string) bool {
    *this = PrivacyValue(sval)
    return true
}

////////////////////////////////////////////////////////////////////////////////

type InputRecord struct {
    Channel string `field:"Channel" required:"false"`
    Video string `field:"Video" required:"false"`
    Stream string `field:"Stream" required:"false"`
    Title string `field:"Title" required:"false"`
    Description string `field:"Description" required:"false"`
    Category string `field:"Category" required:"false"`
    Keywords string `field:"Keywords" required:"false"`
    StartTime string `field:"Start Time" required:"false"`
    Duration string `field:"Duration" required:"false"`
    Thumbnail string `field:"Thumbnail" required:"false"`
    CustomID string `field:"Custom ID" required:"false"`
    Privacy PrivacyValue `field:"Privacy" required:"false"`
    Embedding bool `field:"Embedding" required:"false"`
}

////////////////////////////////////////////////////////////////////////////////

var (
    filename = flag.String("input","","Input filename")
)

////////////////////////////////////////////////////////////////////////////////


func main() {
    // Read flags
    flag.Parse()

    // Make CSV Reader
    fileHandle,err := os.Open(*filename)
    if err != nil {
        log.Fatalf("Error: %v",err)
    }
    defer fileHandle.Close()
    reader := csv.NewReader(fileHandle)

    // make a pointer to our structure
    record := new(InputRecord)
    iterator, err := csvdata.NewReadIter(reader,record)
    if err != nil {
        log.Fatalf("Error: %v",err)
    }
    for iterator.Get() {
        fmt.Println(record)
    }
    if iterator.Error != nil {
        log.Fatalf("Error: %v at line %d and column %d",iterator.Error,iterator.Line,iterator.Column)
    }
}
