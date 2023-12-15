package converter

import (
	"encoding/json"
	"fmt"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/user"
)

func UserToBinary(u *user.User) ([]byte, error) {
	sData := map[string]interface{}{
		"id":            u.ID(),
		"playlists":     u.Playlists(),
		"subscriptions": u.Subscriptions(),
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

	subInterface, ok := uData["subscriptions"]
	if !ok {
		return user.User{}, fmt.Errorf("missing Subscriptions field")
	}
	var subscriptions []string
	if subSlice, ok := subInterface.([]interface{}); ok {
		for _, item := range subSlice {
			if str, isString := item.(string); isString {
				subscriptions = append(subscriptions, str)
			} else {
				return user.User{}, fmt.Errorf("invalid Subscriptions field")
			}
		}
	}

	return user.NewUser(id, playlists, subscriptions), nil
}
