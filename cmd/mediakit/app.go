package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/djthorpe/gopi/v3"
)

type walkfunc func(path string, info os.FileInfo) error

type app struct {
	gopi.Unit
	gopi.MediaManager
	gopi.Logger
	gopi.Command

	offset, limit *uint
}

func (this *app) Define(cfg gopi.Config) error {
	// Set offset and limit
	this.offset = cfg.FlagUint("offset", 0, "File process offset")
	this.limit = cfg.FlagUint("limit", 0, "File process limit")

	// Define commands
	cfg.Command("streams", "Dump stream information", this.Streams)

	return nil
}

func (this *app) New(cfg gopi.Config) error {
	// Set the command
	if this.Command = cfg.GetCommand(nil); this.Command == nil {
		return gopi.ErrHelp
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
			result = append(result, filepath.Clean(arg))
		}
	}
	return result, nil
}

// Walk will traverse through files but only process those within offset/limit
// bounds
func Walk(ctx context.Context, paths []string, count, offset, limit *uint, fn walkfunc) error {
	//seen := make(map[string]bool)

	// Walk through the files
	for _, path := range paths {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			return WalkFunc(ctx, count, offset, limit, path, info, fn, err)
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
	*count += 1
	if *count > *offset {
		if err := fn(path, info); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

/*
func FormatFlags(flags gopi.MediaFlag) string {
	flags_ := fmt.Sprint(flags)
	return strings.Replace(flags_, "|", " ", -1)
}

func FormatArtists(artists []string) string {
	str := ""
	for i, artist := range artists {
		if i > 0 {
			str += ", "
		}
		str += strconv.Quote(artist)
	}
	return str
}

func FormatMetadata(metadata map[gopi.MediaKey]interface{}) string {
	str := ""
	for k, v := range metadata {
		switch v.(type) {
		case string:
			str += fmt.Sprintf(" %s=%q", k, v.(string))
		default:
			str += fmt.Sprintf(" %s=%v", k, v)
		}
	}
	return strings.TrimSpace(str)
}

func (this *album) Artists() []string {
	keys := make(map[string]bool, len(this.files))
	for _, file := range this.files {
		key, ok := file.Metadata[gopi.MEDIA_KEY_ALBUM_ARTIST].(string)
		if ok == false || len(key) < 2 {
			continue
		}
		keys[key] = true
	}
	artists := []string{}
	for k := range keys {
		artists = append(artists, k)
	}
	return artists
}

func (this *file) Filename() string {
	disc, _ := this.Metadata[gopi.MEDIA_KEY_DISC].(uint)
	track, _ := this.Metadata[gopi.MEDIA_KEY_TRACK].(uint)
	title, _ := this.Metadata[gopi.MEDIA_KEY_TITLE].(string)
	ext := filepath.Ext(this.Name)
	str := CleanName(title) + ext
	if track > 0 {
		str = fmt.Sprintf("%02d - %v", track, str)
	}
	if disc > 0 {
		str = fmt.Sprintf("%02d - %v", disc, str)
	}
	return str
}

func CleanName(value string) string {
	value = strings.Replace(value, "/", "_", -1)
	value = strings.Replace(value, ".", "_", -1)
	value = strings.Replace(value, ":", "_", -1)
	return value
}
*/
