package youtubeapi

import "regexp"

func (this *YouTubeService) SetChannel(value string) error {
	// check for empty value
	if value == "" {
		this.channel = value
		return nil
	}
	// check regular expression for channel ID
	matched, _ := regexp.MatchString("^UC([a-z][A-Z][0-9]{22})$", value)
	if matched != true {
		return ErrorBadParameter
	}
	// set parameter
	this.channel = value
	return nil
}
