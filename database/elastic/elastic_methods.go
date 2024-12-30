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

	request := User{
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
