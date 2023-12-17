package command

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"wckd1/tg-youtube-podcasts-bot/internal/converter"
	commandparser "wckd1/tg-youtube-podcasts-bot/internal/delivery/command_parser"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"
	"wckd1/tg-youtube-podcasts-bot/utils"
)

var _ telegram.Command = (*playlist)(nil)

type playlist struct {
	playlistUsecase *usecase.PlaylistUsecase
}

func NewPlaylistCommand(playlistUsecase *usecase.PlaylistUsecase) telegram.Command {
	return playlist{playlistUsecase}
}

// OnMessage return new subscription status
func (p playlist) OnMessage(msg telegram.Message) telegram.Response {
	if !utils.Contains(p.ReactOn(), msg.Command) {
		return telegram.Response{}
	}

	userID := strconv.Itoa(int(msg.ChatID))

	// Parse command arguments
	args, err := commandparser.ParsePlaylistArguments(msg.Arguments)
	if err != nil {
		log.Printf("[ERROR] can't parse playlist arguments, %+v", err)
		return telegram.Response{
			Text: fmt.Sprintf("Can't execute command: %s", err.Error()),
			Send: true,
		}
	}

	// List all playlists
	if len(args) == 0 {
		pls, err := p.playlistUsecase.ListPlaylists(userID)
		if err != nil {
			log.Printf("[ERROR] can't list playlists, %+v", err)
			return telegram.Response{
				Text: "Something went wrong",
				Send: true,
			}
		}

		plsStrings := make([]string, 0)
		for _, pl := range pls {
			plsStrings = append(plsStrings, converter.PlaylistToString(&pl))
		}

		return telegram.Response{
			Text: strings.Join(plsStrings, "\n"),
			Send: true,
		}
	}

	// Create new playlist
	name, ok := args[commandparser.PlaylistNameKey]
	if ok {
		pl, err := p.playlistUsecase.CreatePlaylist(userID, name)
		if err != nil {
			log.Printf("[ERROR] can't create playlists, %+v", err)
			return telegram.Response{
				Text: "Something went wrong",
				Send: true,
			}
		}

		return telegram.Response{
			Text: "Playlist created:\n" + converter.PlaylistToString(pl),
			Send: true,
		}
	}

	// Invalid command
	log.Printf("[ERROR] can't parse playlist arguments, %+v", err)
	return telegram.Response{
		Text: commandparser.ErrInvalidCommand.Error(),
		Send: true,
	}
}

// ReactOn keys
func (p playlist) ReactOn() []string {
	return []string{"pl"}
}
