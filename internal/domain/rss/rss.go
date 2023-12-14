package rss

import (
	"encoding/xml"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
)

type RSS struct {
	XMLName        xml.Name          `xml:"rss"`
	Version        string            `xml:"version,attr"`
	NsItunes       string            `xml:"xmlns:itunes,attr"`
	Title          string            `xml:"channel>title"`
	Description    string            `xml:"channel>description"`
	Image          string            `xml:"channel>itunes:image"`
	Language       string            `xml:"channel>language"`
	ItunesExplicit string            `xml:"channel>itunes:explicit"`
	Category       string            `xml:"channel>itunes:category"`
	ItemList       []episode.Episode `xml:"channel>item"`
}

// TODO: Get rid of xml attributes
// TODO: Add init
