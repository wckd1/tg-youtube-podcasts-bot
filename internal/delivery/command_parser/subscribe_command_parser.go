package commandparser

import (
	"net/url"
	"strings"
)

const (
	SubIDKey       = "id"
	SubURLKey      = "url"
	SubPlaylistKey = "playlist"
)

type SubscribeArguments map[string]string

func ParseSubscribeArguments(arguments string) (AddArguments, error) {
	subArgs := make(AddArguments)
	args := strings.Fields(strings.TrimSpace(arguments))

	// No arguments - invalid command
	if len(args) < 2 {
		return subArgs, ErrInvalidCommand
	}

	// Get url
	id, url, err := validateSubscriptionURL(args[0])
	if err != nil {
		return subArgs, err
	}
	subArgs[SubIDKey] = id
	subArgs[SubURLKey] = url

	// Get playlist if provided
	var playlist string
	for i := 1; i < len(args); i++ {
		if args[i] == "-p" && i+1 < len(args) {
			playlist = args[i+1]
			break
		}
	}
	if playlist != "" {
		subArgs[SubPlaylistKey] = playlist
	}

	return subArgs, nil
}

// Supported links:
//
// Channel:
// /c/{id}
// /channel/{id}
// /channel/@{id}
//
// Playlist:
// /watch?v={video_id}&list={id}
// /playlist?list={id}
func validateSubscriptionURL(urlStr string) (string, string, error) {
	// Parse string to URL
	pURL, err := url.Parse(urlStr)
	if err != nil {
		return "", "", ErrParseURL
	}

	// Check if YouTube link
	if strings.ReplaceAll(pURL.Host, "www.", "") != "youtube.com" {
		return "", "", ErrNotYoutubeURL
	}

	// Check link type
	path := strings.Split(pURL.Path, "/")[1]

	// Check for valid playlist link
	listID, ok := pURL.Query()["list"]
	if ok && (path == "watch" || path == "playlist") {
		id := listID[0]
		url := "https://www.youtube.com/playlist?list=" + listID[0]
		return id, url, nil
	}

	// Check for valid channel link
	if path == "c" || path == "channel" {
		id := strings.Split(pURL.Path, "/")[2]
		url := pURL.String()
		return id, url, nil
	}

	if strings.HasPrefix(path, "@") {
		id := path
		url := pURL.String()
		return id, url, nil
	}

	// No supported links found
	return "", "", ErrNotSupportedURL
}
