package dgraph

import (
	"context"
	"encoding/json"
	"fmt"

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
	if user.Registered == 0 {
		return nil, e.ErrNotFound
	}

	if uid == userUid {
		return user, nil
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
		user_ := edges["users"][0]
		user.IsFollowed = len(user_["is_followed"]) > 0
		user.IsFollowing = len(user_["is_following"]) > 0
	}

	return user, nil
}

func (d *Dgraph) GetUserFollows(ctx context.Context, userUid, after string, reverse bool) ([]string, error) {
	if after == "" {
		after = "0x0"
	}
	vars := map[string]string{
		"$user_uid": userUid,
		"$after":    after,
	}
	var query string
	if reverse {
		query = fmt.Sprintf(userFollowsQuery, "~")
	} else {
		query = fmt.Sprintf(userFollowsQuery, "")
	}
	bytes, err := d.get(ctx, query, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.User
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var userUids []string
	for _, user := range response["users"] {
		for _, follower := range user["followers"] {
			userUids = append(userUids, follower.Uid)
		}
	}

	return userUids, nil
}
