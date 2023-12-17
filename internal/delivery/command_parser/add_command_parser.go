package commandparser

import (
	"errors"
	"net/url"
	"strings"
)

const (
	AddIDKey       = "id"
	AddURLKey      = "url"
	AddPlaylistKey = "playlist"
)

var (
	ErrNoURL           = errors.New("no source url provided")
	ErrParseURL        = errors.New("can't parse URL")
	ErrNotYoutubeURL   = errors.New("only youtube links are supported")
	ErrNotSupportedURL = errors.New("unrecognized link type")
)

type AddArguments map[string]string

func ParseAddArguments(arguments string) (AddArguments, error) {
	addArgs := make(AddArguments)
	args := strings.Fields(strings.TrimSpace(arguments))

	// No arguments - invalid command
	if len(args) < 2 {
		return addArgs, ErrInvalidCommand
	}

	// Get url
	id, err := validateVideoURL(args[0])
	if err != nil {
		return addArgs, err
	}
	addArgs[AddIDKey] = id
	addArgs[AddURLKey] = args[0]

	// Get playlist if provided
	var playlist string
	for i := 1; i < len(args); i++ {
		if args[i] == "-p" && i+1 < len(args) {
			playlist = args[i+1]
			break
		}
	}
	if playlist != "" {
		addArgs[AddPlaylistKey] = playlist
	}

	return addArgs, nil
}

// Supported links:
//
// Video:
// /watch?v={id}
func validateVideoURL(urlStr string) (string, error) {
	// Parse string to URL
	pURL, err := url.Parse(urlStr)
	if err != nil {
		return "", ErrParseURL
	}

	// Check if YouTube link
	if strings.ReplaceAll(pURL.Host, "www.", "") != "youtube.com" {
		return "", ErrNotYoutubeURL
	}

	// Check link type
	path := strings.Split(pURL.Path, "/")[1]

	episodeID := pURL.Query().Get("v")

	if episodeID != "" && path == "watch" {
		return episodeID, nil
	}

	// No supported links found
	return "", ErrNotSupportedURL
}
