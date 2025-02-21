package meili

import (
	"fmt"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/meilisearch/meilisearch-go"
)

func (m *Meili) AddUser(user *models.User) error {
	_, err := m.users.AddDocuments([]map[string]any{{
		"id":       user.Uid,
		"name":     user.Name,
		"username": format(user.Username),
		"hits":     0,
	}})
	return err
}

func (m *Meili) SearchUser(query string) ([]string, error) {
	request := &meilisearch.SearchRequest{Limit: 20}
	results, err := m.users.Search(query, request)
	if err != nil {
		return nil, err
	}
	userUids := make([]string, len(results.Hits))
	for i := range results.Hits {
		hit := results.Hits[i].(map[string]any)
		userUids[i] = hit["id"].(string)
	}
	return userUids, nil
}

func (m *Meili) GetUidByUsername(username string) (string, error) {
	request := &meilisearch.SearchRequest{
		Filter: [][]string{{
			fmt.Sprintf("username = %q", format(username)),
		}},
	}

	results, err := m.users.Search("", request)
	if err != nil {
		return "", err
	}

	if len(results.Hits) == 0 {
		return "", e.ErrNotFound
	} else if len(results.Hits) == 1 {
		hit := results.Hits[0].(map[string]any)
		return hit["id"].(string), nil
	} else {
		return "", e.ErrInvalidMatch
	}
}

func (m *Meili) UpdateUser(user *models.User) error {
	request := &meilisearch.DocumentQuery{}
	var doc map[string]any
	err := m.users.GetDocument(user.Uid, request, &doc)
	if err != nil {
		return err
	}
	doUpdate := false
	if doc["name"] != user.Name {
		doc["name"] = user.Name
		doUpdate = true
	}
	user.Username = format(user.Username)
	if doc["username"] != user.Username {
		doc["username"] = user.Username
		doUpdate = true
	}
	if doUpdate {
		_, err = m.users.UpdateDocuments(doc)
	}
	return err
}

func (m *Meili) HitUser(userUid string) error {
	request := &meilisearch.DocumentQuery{}
	var doc map[string]any
	err := m.users.GetDocument(userUid, request, &doc)
	if err != nil {
		return err
	}
	doc["hits"] = doc["hits"].(float64) + 1
	_, err = m.users.UpdateDocuments(doc)
	return err
}

func (m *Meili) DeleteName(userUid string) error {
	request := &meilisearch.DocumentQuery{}
	var doc map[string]any
	err := m.users.GetDocument(userUid, request, &doc)
	if err != nil {
		return err
	}
	doc["name"] = ""
	_, err = m.users.UpdateDocuments(doc)
	return err
}
