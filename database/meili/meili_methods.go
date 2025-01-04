package meili

import (
	"github.com/dro14/nuqta-service/models"
	"github.com/meilisearch/meilisearch-go"
)

func (m *Meili) Ping() error {
	_, err := m.client.Health()
	return err
}

func (m *Meili) AddUser(user *models.User) error {
	username := user.Username
	if username != "" {
		username = "@" + username
	}
	documents := []map[string]string{{
		"id":       user.Uid,
		"name":     user.Name,
		"username": username,
	}}
	_, err := m.client.Index("users").AddDocuments(documents)
	return err
}

func (m *Meili) SearchUser(query string) ([]any, error) {
	results, err := m.client.Index("users").Search(
		query,
		&meilisearch.SearchRequest{Limit: 10},
	)
	return results.Hits, err
}

func (m *Meili) DeleteUser(uid string) error {
	_, err := m.client.Index("users").DeleteDocument(uid)
	return err
}
