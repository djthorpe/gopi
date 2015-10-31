package main

import (
    "github.com/djthorpe/gopi/csvdata"
    "flag"
    "log"
    "os"
    "encoding/csv"
    "fmt"
)

////////////////////////////////////////////////////////////////////////////////

type InputRecord struct {
    Channel string `field:"Channel"`
    Video string `field:"Video"`
    Stream string `field:"Stream"`
    Title string `field:"Title"`
    Description string `field:"Description"`
    Category string `field:"Category"`
    Keywords string `field:"Keywords"`
    StartTime string `field:"Start Time"`
    Duration string `field:"Duration"`
    Thumbnail string `field:"Thumbnail"`
    CustomID string `field:"Custom ID"`
    Privacy string `field:"Privacy"`
    Embedding bool `field:"Embedding"`
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
