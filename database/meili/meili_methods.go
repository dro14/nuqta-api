package meili

import (
	"strings"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/meilisearch/meilisearch-go"
)

func (m *Meili) AddUser(user *models.User) error {
	username := strings.ToLower(user.Username)
	if username != "" && username[0] != '@' {
		username = "@" + username
	}
	_, err := m.users.AddDocuments([]map[string]any{{
		"id":       user.Uid,
		"name":     user.Name,
		"username": username,
		"hits":     0,
	}})
	return err
}

func (m *Meili) SearchUser(query string) ([]*models.User, error) {
	request := &meilisearch.SearchRequest{Limit: 20}
	results, err := m.users.Search(query, request)
	if err != nil {
		return nil, err
	}
	users := make([]*models.User, 0, len(results.Hits))
	for i := range results.Hits {
		hit := results.Hits[i].(map[string]any)
		users = append(users, &models.User{
			Uid:      hit["id"].(string),
			Name:     hit["name"].(string),
			Username: hit["username"].(string),
		})
	}
	return users, nil
}

func (m *Meili) GetUidByUsername(username string) (string, error) {
	username = strings.ToLower(username)
	if username != "" && username[0] != '@' {
		username = "@" + username
	}

	request := &meilisearch.SearchRequest{
		Filter: "username = '" + username + "'",
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
	username := strings.ToLower(user.Username)
	if username != "" && username[0] != '@' {
		username = "@" + username
	}
	if doc["username"] != username {
		doc["username"] = username
		doUpdate = true
	}
	if doUpdate {
		_, err = m.users.UpdateDocuments(doc)
		return err
	} else {
		return nil
	}
}

func (m *Meili) IncrementHits(uid string) error {
	request := &meilisearch.DocumentQuery{}
	var doc map[string]any
	err := m.users.GetDocument(uid, request, &doc)
	if err != nil {
		return err
	}
	doc["hits"] = doc["hits"].(float64) + 1
	_, err = m.users.UpdateDocuments(doc)
	return err
}

func (m *Meili) DeleteName(uid string) error {
	request := &meilisearch.DocumentQuery{}
	var doc map[string]any
	err := m.users.GetDocument(uid, request, &doc)
	if err != nil {
		return err
	}
	doc["name"] = ""
	_, err = m.users.UpdateDocuments(doc)
	return err
}
