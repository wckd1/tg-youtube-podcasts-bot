package converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"wckd1/tg-youtube-podcasts-bot/internal/domain/entity"
)

func SubscriptionToBinary(s *entity.Subscription) ([]byte, error) {
	sData := map[string]interface{}{
		"id":           s.ID(),
		"url":          s.URL(),
		"filter":       s.Filter(),
		"last_updated": s.LastUpdated().Format(DateFormat),
	}

	return json.Marshal(sData)
}

func BinaryToSubscription(d []byte) (entity.Subscription, error) {
	var sData map[string]interface{}
	if err := json.Unmarshal(d, &sData); err != nil {
		return entity.Subscription{}, err
	}

	id, ok := sData["id"].(string)
	if !ok {
		return entity.Subscription{}, fmt.Errorf("missing or invalid ID field")
	}
	url, ok := sData["url"].(string)
	if !ok {
		return entity.Subscription{}, fmt.Errorf("missing or invalid URL field")
	}
	filter, ok := sData["filter"].(string)
	if !ok {
		return entity.Subscription{}, fmt.Errorf("missing or invalid Filter field")
	}
	lastUpdatedStr, ok := sData["last_updated"].(string)
	if !ok {
		return entity.Subscription{}, fmt.Errorf("missing or invalid Last Updated field")
	}
	lastUpdated, err := time.Parse(DateFormat, lastUpdatedStr)
	if err != nil {
		return entity.Subscription{}, errors.Join(fmt.Errorf("invalid format of Last Updated field"), err)
	}

	return entity.NewSubscription(id, url, filter, lastUpdated), nil
}

func SubscriptionToString(sub *entity.Subscription) string {
	return fmt.Sprintf("id: %s, url: %s, filter: %s", sub.ID(), sub.URL(), sub.Filter())
}
