package youtubeapi

import (
	"regexp"
)

func (this *Parameters) SetMaxResults(value uint) error {
	this.maxresults = value
	return nil
}

func (this *Parameters) SetChannel(value string) error {
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

func (this *Parameters) SetVideo(value string) error {
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

func (this *Parameters) SetStream(value string) error {
	if value != "" {
		// check regular expression for stream key xxxxxxxx-xxxx.xxxx.xxxx.xxxx
		matched, _ := regexp.MatchString("^([a-z0-9])+\\.([a-z0-9]){4}\\-([a-z0-9]){4}\\-([a-z0-9]){4}\\-([a-z0-9]){4}$", value)
		if matched != true {
			return ErrorBadParameter
		}

		// Look up stream key
		var err error
		value, err = this.StreamForKey(value)
		if err != nil {
			return err
		}
	}

	// set parameter
	this.stream = value
	return nil
}

func (this *Parameters) SetStatus(value string) error {
	if value != "" {
		// check regular expression for staus
		matched, _ := regexp.MatchString("^(all|active|completed|upcoming)$", value)
		if matched != true {
			return ErrorBadParameter
		}
	}
	// set parameter
	this.status = value
	return nil
}
