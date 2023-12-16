package converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/episode"
)

func EpisodeToBinary(e *episode.Episode) ([]byte, error) {
	epData := map[string]interface{}{
		"id":           e.ID(),
		"audio_type":   e.AudioType(),
		"url":          e.URL(),
		"link":         e.Link(),
		"cover":        e.Cover(),
		"title":        e.Title(),
		"description":  e.Description(),
		"author":       e.Author(),
		"publish_date": e.PublishDate().Format(DateFormat),
		"length":       strconv.Itoa(e.Length()),
		"duration":     strconv.Itoa(e.Duration()),
	}

	return json.Marshal(epData)
}

func BinaryToEpisode(d []byte) (episode.Episode, error) {
	var epData map[string]interface{}
	if err := json.Unmarshal(d, &epData); err != nil {
		return episode.Episode{}, err
	}

	id, ok := epData["id"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing or invalid ID field")
	}
	audioType, ok := epData["audio_type"].(string)
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

	publishDateStr, ok := epData["publish_date"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing Publish Date field")
	}
	publishDate, err := time.Parse(DateFormat, publishDateStr)
	if err != nil {
		return episode.Episode{}, errors.Join(fmt.Errorf("invalid format of Publish Date field"), err)
	}

	lengthStr, ok := epData["length"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing Length field")
	}
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return episode.Episode{}, fmt.Errorf("invalid Length field, %+v", err)
	}

	durationStr, ok := epData["duration"].(string)
	if !ok {
		return episode.Episode{}, fmt.Errorf("missing Duration field")
	}
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		return episode.Episode{}, fmt.Errorf("invalid Duration field, %+v", err)
	}

	return episode.NewEpisode(id, audioType, url, link, cover, title, description, author, publishDate, length, duration), nil
}
