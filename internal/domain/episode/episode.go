package episode

type Episode struct {
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

// TODO: Get rid of xml attributes
// TODO: Add init
