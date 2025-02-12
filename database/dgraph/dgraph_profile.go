package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
)

var attributes = []string{
	"name",
	"bio",
	"birthday",
	"banner",
	"avatars",
	"thumbnails",
}

func (d *Dgraph) CreateProfile(ctx context.Context, user *models.User) (*models.User, error) {
	user.DType = []string{"User"}
	user.Uid = "_:user"
	user.Username = strings.Split(user.Email, "@")[0]
	user.JoinedAt = time.Now().Unix()

	assigned, err := d.setJson(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Uid = assigned.Uids["user"]
	return user, nil
}

func (d *Dgraph) GetProfile(ctx context.Context, firebaseUid string) (*models.User, error) {
	query := fmt.Sprintf(userByFirebaseUidQuery, firebaseUid)
	bytes, err := d.get(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.User
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	if len(response["users"]) == 0 {
		return nil, e.ErrNotFound
	} else if len(response["users"]) > 1 {
		return nil, e.ErrInvalidMatch
	}

	return response["users"][0], nil
}

func (d *Dgraph) UpdateProfile(ctx context.Context, user *models.User) error {
	user.DType = []string{"User"}
	_, err := d.setJson(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dgraph) DeleteProfileAttribute(ctx context.Context, userUid, attribute, value string) error {
	if !slices.Contains(attributes, attribute) {
		return e.ErrUnknownAttribute
	}
	query := fmt.Sprintf(`<%s> <%s> "%s" .`, userUid, attribute, value)
	return d.deleteNquads(ctx, query)
}
