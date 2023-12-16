package subscription

import "time"

type Subscription struct {
	id          string
	url         string
	filter      string
	lastUpdated time.Time
}

func NewSubscription(id, url, filter string, lastUpdated time.Time) Subscription {
	return Subscription{id, url, filter, lastUpdated}
}

func (s Subscription) ID() string       { return s.id }
func (s *Subscription) SetID(id string) { s.id = id }

func (s Subscription) URL() string        { return s.url }
func (s *Subscription) SetURL(url string) { s.url = url }

func (s Subscription) Filter() string           { return s.filter }
func (s *Subscription) SetFilter(filter string) { s.filter = filter }

func (s Subscription) LastUpdated() time.Time         { return s.lastUpdated }
func (s *Subscription) SetLastUpdated(time time.Time) { s.lastUpdated = time }
