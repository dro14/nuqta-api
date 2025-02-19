package dgraph

import (
	"context"
	"encoding/json"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
)

func (d *Dgraph) GetUser(ctx context.Context, uid, userUid string) (*models.User, error) {
	vars := map[string]string{
		"$user_uid": userUid,
	}
	bytes, err := d.get(ctx, userByUidQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.User
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	user := response["users"][0]
	if user.JoinedAt == 0 {
		return nil, e.ErrNotFound
	}

	vars = map[string]string{
		"$uid":      uid,
		"$user_uid": userUid,
	}
	bytes, err = d.get(ctx, userEdgesQuery, vars)
	if err != nil {
		return nil, err
	}

	var edges map[string][]map[string][]any
	err = json.Unmarshal(bytes, &edges)
	if err != nil {
		return nil, err
	}

	if len(edges["users"]) > 0 {
		user.IsFollowing = len(edges["users"][0]["is_following"]) > 0
		user.IsFollowed = len(edges["users"][0]["is_followed"]) > 0
	}

	return user, nil
}
