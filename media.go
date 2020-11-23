package gopi

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// MediaManager for media file management
type MediaManager interface {
	OpenFile(path string) (Media, error) // Open a media file
	Close(Media) error                   // Close a media file
}

type Media interface {
	//	URL() *url.URL
}
