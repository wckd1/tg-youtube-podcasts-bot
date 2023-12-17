package commandparser

import (
	"errors"
	"strings"
)

var (
	ErrNameIsRequired = errors.New("playlist name is required")
	ErrInvalidCommand = errors.New("invalid command")
)

const (
	PlaylistNameKey = "name"
)

type PlaylistArguments map[string]string

func ParsePlaylistArguments(arguments string) (PlaylistArguments, error) {
	plArgs := make(PlaylistArguments)

	// No arguments - list all playlists
	if len(arguments) < 2 {
		return plArgs, nil
	}

	args := strings.Fields(strings.TrimSpace(arguments))

	switch args[0] {
	case "-new":
		// Create new playlist
		if len(args) < 2 {
			return plArgs, ErrNameIsRequired
		}

		plArgs[PlaylistNameKey] = strings.Join(args[1:], " ")
		return plArgs, nil
	default:
		return plArgs, ErrInvalidCommand
	}
}
