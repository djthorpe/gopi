package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type WorkItem struct {
	Root  string
	Path  string
	Event notify.Event
}

////////////////////////////////////////////////////////////////////////////////
// FLAGS

var (
	flagMediaRoot         = flag.String("mediaroot", "", "Media Root Path")
	flagAllowedExtensions = flag.String("ext", ".m4a .m4v .mov .mp3 .mp4", "Allowed File Extensions")
	flagNumberOfWorkers   = flag.Int("workers", 4, "The number of workers")
)

////////////////////////////////////////////////////////////////////////////////

var allowedExtensionsMap map[string]bool
var workQueue *WorkQueue
var db *Database

////////////////////////////////////////////////////////////////////////////////

func flagUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
}

func checkAllowedExtension(ext string) bool {
	if allowedExtensionsMap == nil {
		extensions := strings.Fields(strings.ToLower(*flagAllowedExtensions))
		allowedExtensionsMap = make(map[string]bool, len(extensions))
		for _, e := range extensions {
			allowedExtensionsMap[e] = true
		}
	}
	if _, ok := allowedExtensionsMap[ext]; ok {
		return true
	}
	return false
}

func walkPath(fullpath string, info os.FileInfo, err error) error {
	// get filename and extension
	filename := path.Base(fullpath)
	ext := strings.ToLower(path.Ext(filename))

	// skip hidden paths
	if strings.HasPrefix(filename, ".") && info.IsDir() {
		return filepath.SkipDir
	}

	// ignore other paths
	if info.IsDir() {
		return nil
	}

	// skip hidden files
	if strings.HasPrefix(filename, ".") {
		return nil
	}

	// check for file extension, queue work
	if checkAllowedExtension(ext) {
		workQueue.Push(WorkItem{Root: *flagMediaRoot, Path: fullpath, Event: notify.Create})
	}

	return nil
}

func processWorkItem(db *Database,item WorkItem) error {

	// skip hidden files
	filename := path.Base(item.Path)
	if strings.HasPrefix(filename, ".") {
		return nil
	}

	// get relative path
	relpath, _ := filepath.Rel(item.Root,item.Path)

	// check for file removal
	_, err := os.Stat(item.Path)
	if os.IsNotExist(err) && item.Event == notify.Remove {
		fmt.Println("Removed: ",item.Path," (",relpath,")")
		return nil
	}

	// check for file addition
	if item.Event == notify.Create && checkAllowedExtension(strings.ToLower(path.Ext(filename))) {
		err := db.Insert(MediaItem{ FullPath: item.Path, RootPath: item.Root, RelPath: relpath })
		if err != nil {
			return err
		}
		return nil
	}

	// check for file rename
	if item.Event == notify.Rename && checkAllowedExtension(strings.ToLower(path.Ext(filename))) {
		fmt.Println("Renamed: ",item.Path," (",relpath,")")
		return nil
	}

	// check for file change
	if item.Event == notify.Write {
		fmt.Println("Modified: ",item.Path," (",relpath,")")
		return nil
	}

	fmt.Println("OTHER ", item.Path, item.Event," (",relpath,")")
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Read flags
	flag.Usage = flagUsage
	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Check mediaroot to ensure it's a directory
	if _, err := os.Stat(*flagMediaRoot); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "-mediaroot: %v\n", err)
		os.Exit(1)
	}

	// Check number of workers which will process media files
	if *flagNumberOfWorkers < 1 {
		fmt.Fprintf(os.Stderr, "-workers: Needs to be a positive value\n")
		os.Exit(1)
	}

	// create database
	db, err := NewDatabase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database error: %v\n", err)
		os.Exit(1)
	}
	defer db.Terminate()

	// create queue for processing media files
	workQueue = NewWorkQueue(*flagNumberOfWorkers, func(val interface{}) {
		err := processWorkItem(db,val.(WorkItem))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	})

	// Walk the tree
	fmt.Println("Creating database")
	if err := filepath.Walk(*flagMediaRoot, walkPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// wait until queue is done
	fmt.Println("Serving database")
	workQueue.Wait()

	// create the path watcher
	watcher := make(chan notify.EventInfo, 1)
	if err := notify.Watch(*flagMediaRoot + "/...", watcher, notify.All); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer notify.Stop(watcher)

	// Process changes to the media root
	for {
		event := <-watcher
		err := processWorkItem(db,WorkItem{Root: *flagMediaRoot, Path: event.Path(), Event: event.Event()})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}
