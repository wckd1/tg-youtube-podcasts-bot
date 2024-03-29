package converter

import (
	"encoding/json"
	"fmt"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
)

func PlaylistToBinary(p *entity.Playlist) ([]byte, error) {
	sData := map[string]interface{}{
		"id":            p.ID(),
		"name":          p.Name(),
		"episodes":      p.Episodes(),
		"subscriptions": p.Subscriptions(),
	}

	return json.Marshal(sData)
}

func BinaryToPlaylist(d []byte) (entity.Playlist, error) {
	var plData map[string]interface{}
	if err := json.Unmarshal(d, &plData); err != nil {
		return entity.Playlist{}, err
	}

	id, ok := plData["id"].(string)
	if !ok {
		return entity.Playlist{}, fmt.Errorf("missing or invalid ID field")
	}

	name, ok := plData["name"].(string)
	if !ok {
		return entity.Playlist{}, fmt.Errorf("missing or invalid Name field")
	}

	epInterfase, ok := plData["episodes"]
	if !ok {
		return entity.Playlist{}, fmt.Errorf("missing Episodes field")
	}
	var episodes []string
	if epSlice, ok := epInterfase.([]interface{}); ok {
		for _, item := range epSlice {
			if str, isString := item.(string); isString {
				episodes = append(episodes, str)
			} else {
				return entity.Playlist{}, fmt.Errorf("invalid Episodes field")
			}
		}
	}

	subInterface, ok := plData["subscriptions"]
	if !ok {
		return entity.Playlist{}, fmt.Errorf("missing Subscriptions field")
	}
	var subscriptions []string
	if subSlice, ok := subInterface.([]interface{}); ok {
		for _, item := range subSlice {
			if str, isString := item.(string); isString {
				subscriptions = append(subscriptions, str)
			} else {
				return entity.Playlist{}, fmt.Errorf("invalid Subscriptions field")
			}
		}
	}

	return entity.NewPlaylist(id, name, episodes, subscriptions), nil
}

func PlaylistToString(pl *entity.Playlist) string {
	return fmt.Sprintf("id: %s, name: %s", pl.ID(), pl.Name())
}
