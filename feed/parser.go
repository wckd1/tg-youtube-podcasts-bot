package feed

import (
	"fmt"
	"net/url"
	"strings"
	"time"
	"wckd1/tg-youtube-podcasts-bot/db"

	"mvdan.cc/xurls/v2"
)

//
// Supported links:
//
// Video:
// /watch?v={id}
//
// Channel:
// /c/{id}
// /channel/{id}
//
// Playlist:
// /watch?v={video_id}&list={id}
// /playlist?list={id}
//

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

	// Check for valid playlist link
	listID, ok := purl.Query()["list"]
	if ok && (path == "watch" || path == "playlist") {
		sub.Type = db.Playlist
		sub.YouTubeID = listID[0]
		return
	}

	// Check for valid video link
	videoID, ok := purl.Query()["v"]
	if ok && path == "watch" {
		sub.Type = db.Video
		sub.YouTubeID = videoID[0]
	}

	// Check for valid channel link
	if path == "c" || path == "channel" {
		sub.Type = db.Channel
		sub.YouTubeID = strings.Split(purl.Path, "/")[2]

		// Parse optional filter
		sub.Filter = strings.ReplaceAll(arguments, furl, "")
	}

	// No supported links found
	err = fmt.Errorf("unrecognized link type")
	return
}

func (fs FeedService) parseDate(input string) string {
	d, _ := time.Parse("20060102", input)
	return d.Format("Mon, 2 Jan 2006 15:04:05 GMT")
}
