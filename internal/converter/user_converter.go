package converter

import (
	"encoding/json"
	"fmt"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
)

func UserToBinary(u *entity.User) ([]byte, error) {
	sData := map[string]interface{}{
		"id":                u.ID(),
		"default_playlists": u.DefaultPlaylist(),
		"playlists":         u.Playlists(),
	}

	return json.Marshal(sData)
}

func BinaryToUser(d []byte) (entity.User, error) {
	var uData map[string]interface{}
	if err := json.Unmarshal(d, &uData); err != nil {
		return entity.User{}, err
	}

	id, ok := uData["id"].(string)
	if !ok {
		return entity.User{}, fmt.Errorf("missing or invalid ID field")
	}

	defaultPl, ok := uData["default_playlists"].(string)
	if !ok {
		return entity.User{}, fmt.Errorf("missing or invalid Default Playlist field")
	}

	plInterface, ok := uData["playlists"]
	if !ok {
		return entity.User{}, fmt.Errorf("missing Playlists field")
	}
	var playlists []string
	if plSlice, ok := plInterface.([]interface{}); ok {
		for _, item := range plSlice {
			if str, isString := item.(string); isString {
				playlists = append(playlists, str)
			} else {
				return entity.User{}, fmt.Errorf("invalid Playlists field")
			}
		}
	}

	return entity.NewUser(id, defaultPl, playlists), nil
}
