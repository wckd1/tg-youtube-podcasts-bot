package feed

import (
	"fmt"
	"net/url"
	"strings"
	"wckd1/tg-youtube-podcasts-bot/db"

	"mvdan.cc/xurls/v2"
)

// Parse params to detect subscription details
func (fs FeedService) parseSubscription(arguments string) (sub db.Subscription, err error) {
	sub = db.Subscription{}

	// Check if arguments contains link
	furl := xurls.Relaxed().FindString(arguments)
	if len(furl) == 0 {
		err = fmt.Errorf("no source url provided")
		return
	}

	// Check if YouTube link
	purl, err := url.Parse(furl)
	if strings.ReplaceAll(purl.Host, "www.", "") != "youtube.com" {
		err = fmt.Errorf("only youtube links are supported")
		return
	}

	// Check link type
	path := strings.Split(purl.Path, "/")[1]
	switch path {
	// Check if passed video link
	case "watch":
		// Check if video is in playlist
		// TODO: Add support for https://www.youtube.com/playlist?list= links
		if _, ok := purl.Query()["list"]; ok {
			err = fmt.Errorf("playlist subscription is not implemented yet")
			// sub.Type = db.Playlist
			// sub.YouTubeID = listID[0]
		} else if videoID, ok := purl.Query()["v"]; ok {
			sub.Type = db.Video
			sub.YouTubeID = videoID[0]
		} else {
			err = fmt.Errorf("unrecognized link type")
		}
	// Check if passed channel link
	case "c", "channel":
		err = fmt.Errorf("channel subscription is not implemented yet")
		// sub.Type = db.Channel

		// // Parse title
		// filter := strings.ReplaceAll(arguments, furl, "")
		// sub.Filter = strings.TrimSpace(filter)

		// channelID := strings.Split(purl.Path, "/")[2]
		// sub.YouTubeID = channelID
	// Other links are unsupported
	default:
		err = fmt.Errorf("unrecognized link type")
	}

	return
}
