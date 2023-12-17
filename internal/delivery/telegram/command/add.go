package command

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	commandparser "wckd1/tg-youtube-podcasts-bot/internal/delivery/command_parser"
	"wckd1/tg-youtube-podcasts-bot/internal/delivery/telegram"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/usecase"
	"wckd1/tg-youtube-podcasts-bot/utils"

	"mvdan.cc/xurls/v2"
)

var (
	ErrNoURL           = errors.New("no source url provided")
	ErrParseURL        = errors.New("can't parse URL")
	ErrNotYoutubeURL   = errors.New("only youtube links are supported")
	ErrNotSupportedURL = errors.New("unrecognized link type")
)

var _ telegram.Command = (*add)(nil)

type add struct {
	addUsecase *usecase.AddUsecase
}

func NewAddCommand(addUsecase *usecase.AddUsecase) telegram.Command {
	return add{addUsecase}
}

// OnMessage return new subscription status
func (a add) OnMessage(msg telegram.Message) telegram.Response {
	if !utils.Contains(a.ReactOn(), msg.Command) {
		return telegram.Response{}
	}

	userID := strconv.Itoa(int(msg.ChatID))

	// Parse command arguments
	args, err := commandparser.ParseAddArguments(msg.Arguments)
	if err != nil {
		log.Printf("[ERROR] can't parse playlist arguments, %+v", err)
		return telegram.Response{
			Text: fmt.Sprintf("Can't execute command: %s", err.Error()),
			Send: true,
		}
	}

	id := args[commandparser.AddIDKey]
	url := args[commandparser.AddURLKey]
	playlist, ok := args[commandparser.AddPlaylistKey]
	if !ok {
		playlist = ""
	}

	if err := a.addUsecase.AddEpisode(userID, id, url, playlist); err != nil {
		log.Printf("[ERROR] failed to add episode. %+v", err)
		return telegram.Response{
			Text: "Failed to add episode",
			Send: true,
		}
	}

	return telegram.Response{
		Text: "Episode added",
		Send: true,
	}
}

// ReactOn keys
func (a add) ReactOn() []string {
	return []string{"add"}
}

type AddItem struct {
	isEpisode bool
	id        string
	url       string
	filter    string
}

// Supported links:
//
// Video:
// /watch?v={id}
//
// Channel:
// /c/{id}
// /channel/{id}
// /channel/@{id}
//
// Playlist:
// /watch?v={video_id}&list={id}
// /playlist?list={id}
func (a add) parseArguments(args string) (AddItem, error) {
	item := AddItem{}

	// Check if arguments contains link
	furl := xurls.Relaxed().FindString(args)
	if len(furl) == 0 {
		return item, ErrNoURL
	}
	purl, err := url.Parse(furl)
	if err != nil {
		return item, errors.Join(ErrParseURL, err)
	}

	// Check if YouTube link
	if strings.ReplaceAll(purl.Host, "www.", "") != "youtube.com" {
		return item, ErrNotYoutubeURL
	}

	// Check link type
	path := strings.Split(purl.Path, "/")[1]

	// Check for valid playlist link
	listID, ok := purl.Query()["list"]
	if ok && (path == "watch" || path == "playlist") {
		item.isEpisode = false
		item.url = "https://www.youtube.com/playlist?list=" + listID[0]
		item.id = listID[0]
		return item, nil
	}

	// Check for valid video link
	episodeID := purl.Query().Get("v")
	if episodeID != "" && path == "watch" {
		item.isEpisode = true
		item.url = purl.String()
		item.id = episodeID
		return item, nil
	}

	// Check for valid channel link
	if path == "c" || path == "channel" {
		item.isEpisode = false
		item.url = purl.String()

		// Parse optional filter
		item.filter = strings.ReplaceAll(args, furl, "")

		chanID := strings.Split(purl.Path, "/")[2]
		item.id = strings.Join([]string{chanID, strings.TrimSpace(item.filter)}, "_")
		return item, nil
	}

	if strings.HasPrefix(path, "@") {
		item.isEpisode = false
		item.url = purl.String()

		// Parse optional filter
		item.filter = strings.ReplaceAll(args, furl, "")
		item.id = strings.Join([]string{path, strings.TrimSpace(item.filter)}, "_")
		return item, nil
	}

	// No supported links found
	return item, ErrNotSupportedURL
}
