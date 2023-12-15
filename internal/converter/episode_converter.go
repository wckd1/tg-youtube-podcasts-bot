package converter

import (
	"encoding/json"
	"fmt"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
	"wckd1/tg-youtube-podcasts-bot/internal/infra/content"
)

func BinaryToEpisode(d []byte) (episode.Episode, error) {
	var epData map[string]interface{}
	if err := json.Unmarshal(d, &epData); err != nil {
		return episode.Episode{}, err
	}

	id, ok := epData["id"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid ID field")
	}
	audioType, ok := epData["audioType"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid Audio Type field")
	}
	url, ok := epData["url"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid URL field")
	}
	link, ok := epData["link"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid Link field")
	}
	cover, ok := epData["cover"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid Cover field")
	}
	title, ok := epData["title"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid Title field")
	}
	description, ok := epData["description"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid Description field")
	}
	author, ok := epData["author"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid Author field")
	}
	publishDate, ok := epData["publishDate"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid Publish Date field")
	}
	length, ok := epData["length"].(int)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid Length field")
	}
	duration, ok := epData["duration"].(int)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid Duration field")
	}

	return episode.NewEpisode(id, audioType, url, link, cover, title, description, author, publishDate, length, duration), nil
}

func DownloadToEpisode(id string, dl content.Download) (episode.Episode, error) {
	pubDate, err := time.Parse("20060102", dl.Info.Date)
	if err != nil {
		return episode.Episode{}, err
	}

	return episode.NewEpisode(
		id,
		dl.Info.Type,
		dl.URL,
		dl.Info.Link,
		dl.Info.ImageURL,
		dl.Info.Title,
		dl.Info.Description, // TODO: "<![CDATA[" + dl.Info.Description + "]]>",
		dl.Info.Author,
		pubDate.Format("Mon, 2 Jan 2006 15:04:05 GMT"), // TODO: Save if DateFormat, parse only on RSS build
		dl.Info.Length,
		dl.Info.Duration,
	), nil
}
