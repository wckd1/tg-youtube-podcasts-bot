package converter

import (
	"encoding/json"
	"fmt"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
)

func UserToBinary(u *user.User) ([]byte, error) {
	sData := map[string]interface{}{
		"id":                u.ID(),
		"default_playlists": u.DefaultPlaylist(),
		"playlists":         u.Playlists(),
	}

	return json.Marshal(sData)
}

func BinaryToUser(d []byte) (user.User, error) {
	var uData map[string]interface{}
	if err := json.Unmarshal(d, &uData); err != nil {
		return user.User{}, err
	}

	id, ok := uData["id"].(string)
	if !ok {
		return user.User{}, fmt.Errorf("missing or invalid ID field")
	}

	defaultPl, ok := uData["default_playlists"].(string)
	if !ok {
		return user.User{}, fmt.Errorf("missing or invalid Default Playlist field")
	}

	plInterface, ok := uData["playlists"]
	if !ok {
		return user.User{}, fmt.Errorf("missing Playlists field")
	}
	var playlists []string
	if plSlice, ok := plInterface.([]interface{}); ok {
		for _, item := range plSlice {
			if str, isString := item.(string); isString {
				playlists = append(playlists, str)
			} else {
				return user.User{}, fmt.Errorf("invalid Playlists field")
			}
		}
	}

	return user.NewUser(id, defaultPl, playlists), nil
}
