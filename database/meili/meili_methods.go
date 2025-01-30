package meili

import (
	"github.com/dro14/nuqta-service/models"
	"github.com/meilisearch/meilisearch-go"
)

func (m *Meili) AddUser(user *models.User) error {
	username := user.Username
	if username != "" {
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
	if doc["username"] != user.Username {
		doc["username"] = user.Username
		doUpdate = true
	}
	if doUpdate {
		_, err = m.users.UpdateDocuments(doc)
	}
	return err
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

func (m *Meili) DeleteUser(uid string) error {
	_, err := m.users.DeleteDocument(uid)
	return err
}
