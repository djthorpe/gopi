package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/djthorpe/gopi/v3"
)

type walkfunc func(path string, info os.FileInfo) error

type app struct {
	gopi.Unit
	gopi.MediaManager
	gopi.Logger
	gopi.Command

	offset, limit *uint   // File processing offsets
	quiet         *bool   // Whether errors should be displayed
	match         *string // Regular expression to match

	regexp *regexp.Regexp // Regular expression for filename
	fields *fields        // Metadata which should be displayed
}

func (this *app) Define(cfg gopi.Config) error {
	// Set command-line flags
	this.offset = cfg.FlagUint("offset", 0, "File process offset")
	this.limit = cfg.FlagUint("limit", 0, "File process limit")
	this.quiet = cfg.FlagBool("quiet", false, "Don't display file scan errors")
	this.match = cfg.FlagString("match", "", "Match filenames regular expression")

	// Define commands
	cfg.Command("metadata", "Dump metadata information", this.Metadata)
	cfg.Command("streams", "Dump stream information", this.Streams)
	cfg.Command("thumbnails", "Extract thumbnails", this.Thumbnails)

	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	// Set the command
	if this.Command = cfg.GetCommand(nil); this.Command == nil {
		return gopi.ErrHelp
	}

	// Set up fields
	this.fields = NewFields()

	// Set up regular expression for matching filenames
	if *this.match != "" {
		if re, err := regexp.Compile(*this.match); err != nil {
			return err
		} else {
			this.regexp = re
		}
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

// GetFileArgs returns all files in arguments, or current
// working directory if no arguments provided
func GetFileArgs(args []string) ([]string, error) {
	// Default to the current working directory
	if cwd, err := os.Getwd(); err != nil {
		return nil, err
	} else if len(args) == 0 {
		return []string{cwd}, nil
	}

	// Append files and folders, normalizing them to absolute paths
	result := make([]string, 0, len(args))
	for _, arg := range args {
		if _, err := os.Stat(arg); os.IsNotExist(err) {
			return nil, fmt.Errorf("%q: %w", filepath.Base(arg), gopi.ErrNotFound)
		} else if err != nil {
			return nil, fmt.Errorf("%q: %w", filepath.Base(arg), err)
		} else if filepath.IsAbs(arg) == false {
			if abs, err := filepath.Abs(arg); err == nil {
				arg = abs
			}
		}
		result = append(result, filepath.Clean(arg))
	}
	return result, nil
}

// Walk will traverse through files but only process those within offset/limit
// bounds
func (this *app) Walk(ctx context.Context, paths []string, count *uint, fn walkfunc) error {
	// Walk through the files
	for _, path := range paths {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if this.regexp == nil || this.regexp.MatchString(path) {
				return WalkFunc(ctx, count, this.offset, this.limit, path, info, fn, err)
			} else {
				return nil
			}
		}); err != nil && err != io.EOF {
			return err
		}
	}

	// Return success
	return nil
}

func WalkFunc(ctx context.Context, count, offset, limit *uint, path string, info os.FileInfo, fn walkfunc, err error) error {
	// Deal with incoming errors
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		return err
	}

	// Ignore hidden files and folders
	if strings.HasPrefix(info.Name(), ".") {
		if info.IsDir() {
			return filepath.SkipDir
		} else {
			return nil
		}
	}

	// Ignore anything which isn't a regular file
	if info.Mode().IsRegular() == false {
		return nil
	}

	// If limit has been reached, return io.EOF
	if *limit > 0 && *count >= *limit {
		return io.EOF
	}

	// Increment the count and check
	if *count++; *count > *offset {
		if err := fn(path, info); err != nil {
			return err
		}
	}

	// Return success
	return nil
}
