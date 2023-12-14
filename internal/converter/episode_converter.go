package converter

import (
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/content"
)

func DownloadToEpisode(dl content.Download) (episode.Episode, error) {
	PubDate, err := time.Parse(DateFormat, dl.Info.Date)
	if err != nil {
		return episode.Episode{}, err
	}

	return episode.Episode{
		Enclosure: episode.Enclosure{
			URL:    dl.URL,
			Length: dl.Info.Length,
			Type:   dl.Info.Type,
		},
		Link:        dl.Info.Link,
		Image:       dl.Info.ImageURL,
		Title:       dl.Info.Title,
		Description: "<![CDATA[" + dl.Info.Description + "]]>",
		Author:      dl.Info.Author,
		Duration:    dl.Info.Duration,
		PubDate:     PubDate.Format("Mon, 2 Jan 2006 15:04:05 GMT"), // TODO: Save if DateFormat, parse only on RSS build
	}, nil
}
