package dgraph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
)

func (d *Dgraph) GetUserByUid(ctx context.Context, firebaseUid, uid string) (*models.User, error) {
	query := fmt.Sprintf(userByUidQuery, uid)
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.User
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	user := response["users"][0]
	if user.JoinedAt == 0 {
		return nil, e.ErrNotFound
	}

	user.IsFollowed, err = d.DoesEdgeExist(ctx, uid, "~follow", firebaseUid)
	if err != nil {
		return nil, err
	}

	return user, nil
}
