package main

import (
    "flag"
	"os"
	"log"
	"fmt"
	"path"
	"strings"
	"path/filepath"
	"gopkg.in/fsnotify.v1"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type WorkItem struct {
	Path string
}

////////////////////////////////////////////////////////////////////////////////
// FLAGS

var (
    flagMediaRoot = flag.String("mediaroot","","Media Root Path")
	flagAllowedExtensions = flag.String("ext",".m4a .m4v .mov .mp3 .mp4","Allowed File Extensions")
	flagNumberOfWorkers = flag.Int("workers",4,"The number of workers")
)

////////////////////////////////////////////////////////////////////////////////

var allowedExtensionsMap map[string]bool
var workQueue *WorkQueue

////////////////////////////////////////////////////////////////////////////////

func flagUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n",path.Base(os.Args[0]))
	flag.PrintDefaults()
}

func checkAllowedExtension(ext string) bool {
	if allowedExtensionsMap == nil {
		extensions := strings.Fields(strings.ToLower(*flagAllowedExtensions))
		allowedExtensionsMap = make(map[string]bool,len(extensions))
		for _,e := range extensions {
			allowedExtensionsMap[e] = true
		}
	}
	if _,ok := allowedExtensionsMap[ext]; ok {
		return true
	}
	return false
}

func walkPath(fullpath string,info os.FileInfo,err error) error {
	// get filename and extension
	filename := path.Base(fullpath)
	ext := strings.ToLower(path.Ext(filename))

	// skip hidden paths
	if strings.HasPrefix(filename,".") && info.IsDir() {
		return filepath.SkipDir
	}

	// ignore other paths
	if info.IsDir() {
		return nil
	}

	// skip hidden files
	if strings.HasPrefix(filename,".") {
		return nil
	}

	// check for file extension, queue work
	if checkAllowedExtension(ext) {
		workQueue.Push(WorkItem{Path:fullpath})
	}

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
		fmt.Fprintf(os.Stderr,"-mediaroot: %v\n",err)
		os.Exit(1)
	}

	// Check number of workers
	if *flagNumberOfWorkers < 1 {
		fmt.Fprintf(os.Stderr,"-workers: Needs to be a positive value\n")
		os.Exit(1)
	}

	// create queue, set up work handler
	workQueue = NewWorkQueue(*flagNumberOfWorkers,func(val interface{}) {
		processWorkItem(val.(WorkItem))
	})

	// Walk the tree
    if err := filepath.Walk(*flagMediaRoot,walkPath); err != nil {
        log.Fatal(err)
    }

	// wait until queue is done
	workQueue.Wait()
}
