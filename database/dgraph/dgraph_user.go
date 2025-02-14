package dgraph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
)

func (d *Dgraph) GetUser(ctx context.Context, uid, userUid string) (*models.User, error) {
	query := fmt.Sprintf(userByUidQuery, userUid)
	bytes, err := d.get(ctx, query)
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

	user.IsFollowed, err = d.GetEdge(ctx, uid, "follow", userUid)
	if err != nil {
		return nil, err
	}

	user.IsFollowing, err = d.GetEdge(ctx, userUid, "follow", uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}
