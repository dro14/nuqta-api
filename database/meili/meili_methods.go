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
	documents := []map[string]any{{
		"id":       user.Uid,
		"name":     user.Name,
		"username": username,
		"hits":     0,
	}}
	_, err := m.index.AddDocuments(documents)
	return err
}

func (m *Meili) SearchUser(query string) ([]*models.User, error) {
	results, err := m.index.Search(
		query,
		&meilisearch.SearchRequest{Limit: 10},
	)
	var users []*models.User
	for i := range results.Hits {
		hit := results.Hits[i].(map[string]any)
		users = append(users, &models.User{
			Uid:      hit["id"].(string),
			Name:     hit["name"].(string),
			Username: hit["username"].(string),
		})
	}
	return users, err
}

func (m *Meili) IncrementUserHits(uid string) error {
	var doc map[string]any
	err := m.index.GetDocument(
		uid,
		&meilisearch.DocumentQuery{},
		&doc,
	)
	if err != nil {
		return err
	}
	doc["hits"] = doc["hits"].(float64) + 1
	_, err = m.index.UpdateDocuments(
		[]map[string]any{doc},
		uid,
	)
	return err
}

func (m *Meili) DeleteUser(uid string) error {
	_, err := m.index.DeleteDocument(uid)
	return err
}
