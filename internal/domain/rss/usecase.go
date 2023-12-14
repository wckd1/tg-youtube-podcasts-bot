package rss

import (
	"encoding/xml"
	"errors"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
)

var ErrXMLEncoding = errors.New("failed to build xml")

type RSSUseCase struct{}

func NewRSSUseCase() *RSSUseCase {
	return &RSSUseCase{}
}

func (uc RSSUseCase) BuildRSS() (string, error) {
	// TODO: Get actual items
	items := make([]episode.Episode, 0)

	rss := RSS{
		Version:        "2.0",
		NsItunes:       "http://www.itunes.com/dtds/podcast-1.0.dtd",
		Title:          "Private RSS feed",
		Description:    "Generated by tg-youtube-podcasts-bot",
		Image:          "https://www.clipartkey.com/mpngs/m/197-1971515_youtube-music-seamless-audio-video-switching-transparent-youtube.png",
		Language:       "ru",
		ItunesExplicit: "false",
		Category:       "Education",
		ItemList:       items,
	}

	b, err := xml.MarshalIndent(&rss, "", "  ")
	if err != nil {
		return "", errors.Join(ErrXMLEncoding, err)
	}

	res := `<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(b)
	return res, nil
}
