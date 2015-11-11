package youtubeapi

import (
	"regexp"
)

func (this *YouTubeService) SetMaxResults(value uint) *YouTubeService {
	this.maxresults = value
	return this
}

func (this *YouTubeService) SetChannel(value string) error {
	if value != "" {
		// check regular expression for channel ID
		matched, _ := regexp.MatchString("^UC([a-zA-Z0-9]){22}$", value)
		if matched != true {
			return ErrorBadParameter
		}
	}
	// set parameter
	this.channel = value
	return nil
}

func (this *YouTubeService) SetVideo(value string) error {
	if value != "" {
		// check regular expression for video ID
		matched, _ := regexp.MatchString("^([a-zA-Z0-9\\_\\-]){11}$", value)
		if matched != true {
			return ErrorBadParameter
		}
	}
	// set parameter
	this.video = value
	return nil
}

func (this *YouTubeService) SetStream(value string) error {
	if value != "" {
		// check regular expression for stream key xxxxxxxx-xxxx.xxxx.xxxx.xxxx
		matched, _ := regexp.MatchString("^([a-z0-9])+\\.([a-z0-9]){4}\\-([a-z0-9]){4}\\-([a-z0-9]){4}\\-([a-z0-9]){4}$", value)
		if matched != true {
			return ErrorBadParameter
		}
	}
	// set parameter
	this.stream = value
	return nil
}

func (this *YouTubeService) SetStatus(value string) error {
	if value != "" {
		// check regular expression for staus
		matched, _ := regexp.MatchString("^(all|active|completed|upcoming)$",value)
		if matched != true {
			return ErrorBadParameter
		}
	}
	// set parameter
	this.status = value
	return nil
}
