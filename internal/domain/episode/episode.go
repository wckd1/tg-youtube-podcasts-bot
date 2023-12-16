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

func (e Episode) ID() string          { return e.id }
func (e Episode) Length() int         { return e.length }
func (e Episode) AudioType() string   { return e.audioType }
func (e Episode) URL() string         { return e.url }
func (e Episode) Link() string        { return e.link }
func (e Episode) Cover() string       { return e.cover }
func (e Episode) Title() string       { return e.title }
func (e Episode) Description() string { return e.description }
func (e Episode) Author() string      { return e.author }
func (e Episode) Duration() int       { return e.duration }
func (e Episode) PublishDate() string { return e.publishDate }
