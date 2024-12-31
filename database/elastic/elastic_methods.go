package elastic

import (
	"context"
	"net/http"

	"github.com/dro14/nuqta-service/models"
)

func (e *Elastic) CreateUser(ctx context.Context, user *models.User) error {
	name := user.Name
	username := user.Username
	if len(username) > 0 {
		username = "@" + username
	}

	request := Doc{
		Name:     name,
		Username: username,
		HitCount: 0,
	}

	endpoint := "/users/_doc/" + id(ctx)
	_, err := e.makeRequest(ctx, request, endpoint, http.MethodPut)
	if err != nil {
		return err
	}

	return nil
}

func (e *Elastic) SearchUser(ctx context.Context, query string) ([]string, error) {
	request := searchRequest(query)
	endpoint := "/users/_search"

	response, err := e.makeRequest(ctx, request, endpoint, http.MethodGet)
	if err != nil {
		return nil, err
	}

	ids, err := searchResponse(response)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
