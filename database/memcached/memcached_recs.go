package memcached

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dro14/nuqta-service/models"
)

func (m *Memcached) SetRecs(recs []*models.Post) error {
	json, err := json.Marshal(recs)
	if err != nil {
		return err
	}
	item := &memcache.Item{Key: "recs", Value: json}
	return m.client.Set(item)
}

func (m *Memcached) GetRecs() ([]*models.Post, error) {
	item, err := m.client.Get("recs")
	if errors.Is(err, memcache.ErrCacheMiss) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var recs []*models.Post
	err = json.Unmarshal(item.Value, &recs)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func (m *Memcached) GetRecUpdateTime() (time.Time, error) {
	now := time.Now()
	item, err := m.client.Get("rec_update_time")
	if errors.Is(err, memcache.ErrCacheMiss) {
		value := []byte(strconv.FormatInt(now.Unix(), 10))
		item = &memcache.Item{Key: "rec_update_time", Value: value}
		return now, m.client.Set(item)
	} else if err != nil {
		return now, err
	}
	timestamp, err := strconv.ParseInt(string(item.Value), 10, 64)
	if err != nil {
		return now, err
	}
	return time.Unix(timestamp, 0), nil
}

func (m *Memcached) IncrRecUpdateTime() (time.Time, error) {
	timestamp, err := m.client.Increment("rec_update_time", 60)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(int64(timestamp), 0), nil
}
