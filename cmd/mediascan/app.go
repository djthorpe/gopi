package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/djthorpe/gopi/v3"
	"github.com/olekukonko/tablewriter"
)

type app struct {
	gopi.Unit
	gopi.MediaManager
	gopi.Logger
	gopi.Command

	offset, limit *uint
	files         []*file
	albums        map[string]*album
}

type file struct {
	Path     string
	Name     string
	Error    error
	Flags    gopi.MediaFlag
	Metadata map[gopi.MediaKey]interface{}
	Streams  []string
}

type album struct {
	Name    string
	files   []*file
	artists []string
	folder  string
}

func (this *app) Define(cfg gopi.Config) error {
	// Set offset and limit
	this.offset = cfg.FlagUint("offset", 0, "File process offset")
	this.limit = cfg.FlagUint("limit", 0, "File process limit")

	// Define commands
	cfg.Command("scan", "Scan folder and display information", this.Scan)
	cfg.Command("mkdir", "Make folders for music files", this.Mkdir)

	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if this.Command = cfg.GetCommand(nil); this.Command == nil {
		return gopi.ErrHelp
	}

	// Make maps
	this.albums = make(map[string]*album)

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

func (this *app) Scan(ctx context.Context) error {
	// Set offset counter
	count := uint(0)

	// Walk through the files
	if folder, err := GetFolder(this.Command.Args()); err != nil {
		return err
	} else if err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		return this.WalkFunc(ctx, &count, path, info, err)
	}); err != nil {
		return err
	}

	// Print processed files
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Type", "Metadata", "Streams"})
	for _, file := range this.files {
		if file.Error != nil {
			table.Append([]string{
				strconv.Quote(file.Name),
				fmt.Sprint(file.Error),
				"",
				"",
			})
		} else {
			table.Append([]string{
				strconv.Quote(file.Name),
				FormatFlags(file.Flags),
				FormatMetadata(file.Metadata),
				strings.Join(file.Streams, " "),
			})
		}
	}
	table.Render()
	if len(this.files) == 0 {
		fmt.Fprintln(os.Stdout, "No files processed")
	} else {
		fmt.Fprintf(os.Stdout, "Files %v-%v of %v\n", *this.offset+1, (*this.offset)+uint(len(this.files)), count-1)
	}

	// Return success
	return nil
}

func (this *app) Mkdir(ctx context.Context) error {
	// Set offset counter
	count := uint(0)

	// Walk through the files
	root, err := GetFolder(this.Command.Args())
	if err != nil {
		return err
	} else if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		return this.WalkFunc(ctx, &count, path, info, err)
	}); err != nil {
		return err
	}

	// Process albums
	for _, file := range this.files {
		if file.Flags&gopi.MEDIA_FLAG_ALBUM_TRACK == 0 {
			continue
		}
		key, ok := file.Metadata[gopi.MEDIA_KEY_ALBUM].(string)
		if ok == false || len(key) < 2 || key == "Various Artists" { // Ignore one-character albums
			continue
		} else if _, exists := this.albums[key]; exists == false {
			this.albums[key] = &album{Name: strings.TrimSpace(key)}
		}

		album := this.albums[key]
		album.files = append(album.files, file)
	}

	// Process artists and get folder name
	for _, album := range this.albums {
		album.artists = album.Artists()
		if len(album.artists) == 0 {
			continue
		} else if len(album.artists) == 1 {
			album.folder = filepath.Join(CleanName(album.artists[0]), CleanName(album.Name))
		} else {
			album.folder = filepath.Join("Compilations", CleanName(album.Name))
		}
	}

	// Create folders
	for _, album := range this.albums {
		folder := filepath.Clean(filepath.Join(root, album.folder))
		if err := os.MkdirAll(folder, 0775); err != nil {
			fmt.Fprintln(os.Stderr, folder, err)
		}
	}

	// Move files
	for _, album := range this.albums {
		folder := filepath.Clean(filepath.Join(root, album.folder))
		for _, file := range album.files {
			dest := filepath.Join(folder, file.Filename())
			if err := os.Rename(file.Path, dest); err != nil {
				fmt.Fprintln(os.Stderr, file.Path, err)
			}
		}
	}

	// Print albums, artists and number of tracks
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Album", "Artists", "Folder", "Files"})
	for _, album := range this.albums {
		table.Append([]string{
			strconv.Quote(album.Name),
			FormatArtists(album.artists),
			strconv.Quote(album.folder),
			fmt.Sprint(len(album.files)),
		})
	}
	table.Render()

	// Return success
	return nil
}

func GetFolder(args []string) (string, error) {
	if len(args) == 0 {
		return os.Getwd()
	}
	if len(args) > 1 {
		return "", gopi.ErrHelp
	}
	if _, err := os.Stat(args[0]); err != nil {
		return "", err
	} else {
		return args[0], nil
	}
}

func (this *app) WalkFunc(ctx context.Context, count *uint, path string, info os.FileInfo, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		this.Print("Error: ", info.Name(), err)
		return err
	}
	if strings.HasPrefix(info.Name(), ".") {
		if info.IsDir() {
			return filepath.SkipDir
		} else {
			return nil
		}
	}
	if filepath.Ext(info.Name()) == ".csv" {
		return nil
	}

	// Increment the count and check
	if *count > *this.offset {
		if *this.limit == 0 || uint(len(this.files)) < *this.limit {
			if info.Mode().IsRegular() {
				this.files = append(this.files, this.Read(path, info))
			}
		}
	}

	// Increment the counter
	*count += 1

	// Return success
	return nil
}

func (this *app) Read(path string, info os.FileInfo) *file {
	f := new(file)
	f.Path = path
	f.Name = info.Name()
	f.Metadata = make(map[gopi.MediaKey]interface{}, 10)
	if media, err := this.MediaManager.OpenFile(path); err != nil {
		f.Error = err
	} else {
		defer this.MediaManager.Close(media)
		f.Flags = media.Flags()
		for _, k := range media.Metadata().Keys() {
			if k == "id3v2_priv.www.amazon.com" || k == "iTunSMPB" || k == "iTunNORM" || k == "iTunes_CDDB_1" {
				continue
			}
			f.Metadata[k] = media.Metadata().Value(k)
		}
		for _, s := range media.Streams() {
			f.Streams = append(f.Streams, fmt.Sprint(s))
		}
	}
	return f
}

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
