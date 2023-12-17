package commandparser

import (
	"net/url"
	"strings"
)

const (
	SubIDKey       = "id"
	SubURLKey      = "url"
	SubPlaylistKey = "playlist"
	SubFilterKey   = "filter"
)

type SubscribeArguments map[string]string

func ParseSubscribeArguments(arguments string) (AddArguments, error) {
	subArgs := make(AddArguments)
	args := strings.Fields(strings.TrimSpace(arguments))

	// No arguments - list default subscriptions
	if len(args) < 2 {
		return subArgs, nil
	}

	// Get additional parameters if provided
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-new":
			if i+1 < len(args) {
				id, url, err := validateSubscriptionURL(args[i+1])
				if err != nil {
					return subArgs, err
				}
				subArgs[SubIDKey] = id
				subArgs[SubURLKey] = url
				i++
			}
		case "-p":
			var playlist string
			for j := i + 1; j < len(args) && !strings.HasPrefix(args[j], "-"); j++ {
				playlist += args[j] + " "
			}
			subArgs[SubPlaylistKey] = strings.TrimSpace(playlist)
			i += strings.Count(playlist, " ")
		case "-f":
			var filter string
			for j := i + 1; j < len(args) && !strings.HasPrefix(args[j], "-"); j++ {
				filter += args[j] + " "
			}
			subArgs[SubFilterKey] = strings.TrimSpace(filter)
			i += strings.Count(filter, " ")
		}
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
