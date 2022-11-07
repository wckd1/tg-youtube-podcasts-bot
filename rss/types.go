package rss

import "encoding/xml"

type RSS struct {
	XMLName        xml.Name `xml:"rss"`
	Version        string   `xml:"version,attr"`
	NsItunes       string   `xml:"xmlns:itunes,attr"`
	Title          string   `xml:"channel>title"`
	Description    string   `xml:"channel>description"`
	Link           string   `xml:"channel>link"`
	Language       string   `xml:"channel>language"`
	ItunesExplicit string   `xml:"channel>itunes:explicit"`
	ItemList       []Item   `xml:"channel>item"`
}

type Item struct {
	GUID        string    `xml:"guid"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	Enclosure   Enclosure `xml:"enclosure"`
	// PubDate  string `xml:"pubDate,omitempty"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}
