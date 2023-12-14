package converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/subscription"

	"mvdan.cc/xurls/v2"
)

var (
	ErrNoURL           = errors.New("no source url provided")
	ErrParseURL        = errors.New("can't parse URL")
	ErrNotYoutubeURL   = errors.New("only youtube links are supported")
	ErrNotSupportedURL = errors.New("unrecognized link type")
)

func SubscriptionToBinary(s *subscription.Subscription) ([]byte, error) {
	sData := map[string]interface{}{
		"id":           s.ID(),
		"url":          s.URL(),
		"filter":       s.Filter(),
		"last_updated": s.LastUpdated().Format(DateFormat),
	}

	return json.Marshal(sData)
}

func BinaryToSubscription(d []byte) (subscription.Subscription, error) {
	var sData map[string]interface{}
	if err := json.Unmarshal(d, &sData); err != nil {
		return subscription.Subscription{}, err
	}

	id, ok := sData["id"].(string)
	if !ok {
		return subscription.Subscription{}, fmt.Errorf("missing or invalid ID field")
	}
	url, ok := sData["url"].(string)
	if !ok {
		return subscription.Subscription{}, fmt.Errorf("missing or invalid URL field")
	}
	filter, ok := sData["filter"].(string)
	if !ok {
		return subscription.Subscription{}, fmt.Errorf("missing or invalid Filter field")
	}
	lastUpdatedStr, ok := sData["last_updated"].(string)
	if !ok {
		return subscription.Subscription{}, fmt.Errorf("missing or invalid Last Updated field")
	}
	lastUpdated, err := time.Parse(DateFormat, lastUpdatedStr)
	if err != nil {
		return subscription.Subscription{}, errors.Join(fmt.Errorf("invalid format of Last Updated field"), err)
	}

	return subscription.NewSubscription(id, url, filter, lastUpdated), nil
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
func BotArgumentToSubscription(args string) (subscription.Subscription, error) {
	sub := subscription.Subscription{}

	// Check if arguments contains link
	furl := xurls.Relaxed().FindString(args)
	if len(furl) == 0 {
		return sub, ErrNoURL
	}
	purl, err := url.Parse(furl)
	if err != nil {
		return sub, errors.Join(ErrParseURL, err)
	}

	// Check if YouTube link
	if strings.ReplaceAll(purl.Host, "www.", "") != "youtube.com" {
		return sub, ErrNotYoutubeURL
	}

	// Check link type
	path := strings.Split(purl.Path, "/")[1]

	// Check for valid playlist link
	listID, ok := purl.Query()["list"]
	if ok && (path == "watch" || path == "playlist") {
		sub.SetID(listID[0])
		sub.SetURL("https://www.youtube.com/playlist?list=" + listID[0])
		return sub, nil
	}

	// Check for valid video link
	_, ok = purl.Query()["v"]
	if ok && path == "watch" {
		// sub.IsVideo = true
		sub.SetURL(purl.String())
		return sub, nil
	}

	// Check for valid channel link
	if path == "c" || path == "channel" {
		sub.SetURL(purl.String())

		// Parse optional filter
		filter := strings.ReplaceAll(args, furl, "")
		sub.SetFilter(strings.TrimSpace(filter))

		chanID := strings.Split(purl.Path, "/")[2]
		sub.SetID(strings.Join([]string{chanID, filter}, "_"))
		return sub, nil
	}

	if strings.HasPrefix(path, "@") {
		sub.SetURL(purl.String())

		// Parse optional filter
		filter := strings.ReplaceAll(args, furl, "")
		sub.SetFilter(strings.TrimSpace(filter))

		sub.SetID(strings.Join([]string{path, filter}, "_"))
		return sub, nil
	}

	// No supported links found
	return sub, ErrNotSupportedURL
}
