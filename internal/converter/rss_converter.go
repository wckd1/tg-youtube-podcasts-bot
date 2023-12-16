package converter

import (
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/rss"
)

func EpisodeToRSSEpisode(ep *episode.Episode) rss.RSSEpisode {
	return rss.RSSEpisode{
		UUID: ep.ID(),
		Enclosure: rss.Enclosure{
			URL:    ep.URL(),
			Length: ep.Length(),
			Type:   ep.AudioType(),
		},
		Link:        ep.Link(),
		Image:       ep.Cover(),
		Title:       ep.Title(),
		Description: "<![CDATA[" + ep.Description() + "]]>",
		Author:      ep.Author(),
		Duration:    ep.Duration(),
		PubDate:     ep.PublishDate().Format("Mon, 2 Jan 2006 15:04:05 GMT"),
	}
}
