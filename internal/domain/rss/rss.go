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

type RSSEpisode struct {
	UUID        string    `xml:"guid"`
	Enclosure   Enclosure `xml:"enclosure"`
	Link        string    `xml:"link"`
	Image       string    `xml:"image"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Author      string    `xml:"author,omitempty"`
	Duration    int       `xml:"duration,omitempty"`
	PubDate     string    `xml:"pubDate,omitempty"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}
