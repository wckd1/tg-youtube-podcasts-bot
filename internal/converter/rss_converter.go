package converter

import (
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
)

func EpisodeToRSSEpisode(ep *entity.Episode) entity.RSSEpisode {
	return entity.RSSEpisode{
		UUID: ep.ID(),
		Enclosure: entity.Enclosure{
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
