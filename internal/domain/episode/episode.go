package episode

type Episode struct {
	id          string
	length      int
	audioType   string
	url         string
	link        string
	cover       string
	title       string
	description string
	author      string
	duration    int
	publishDate string
}

func NewEpisode(id, audioType, url, link, cover, title, description, author, publishDate string, length, duration int) Episode {
	return Episode{id, length, audioType, url, link, cover, title, description, author, duration, publishDate}
}

func (e Episode) ID() string { return e.id }
