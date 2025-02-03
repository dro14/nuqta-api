package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
)

var attributes = []string{
	"name",
	"bio",
	"birthday",
	"banner",
	"avatars",
}

func (d *Dgraph) CreateProfile(ctx context.Context, user *models.User) (*models.User, error) {
	user.DType = []string{"User"}
	user.Uid = "_:user"
	json, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	mutation := &api.Mutation{SetJson: json, CommitNow: true}
	assigned, err := d.client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return nil, err
	}

	user.Uid = assigned.Uids["user"]
	return user, nil
}

func (d *Dgraph) GetProfile(ctx context.Context, firebaseUid string) (*models.User, error) {
	query := fmt.Sprintf(userByFirebaseUidQuery, firebaseUid)
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.User
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	if len(response["users"]) == 0 {
		return nil, e.ErrNotFound
	} else if len(response["users"]) == 1 {
		return response["users"][0], nil
	} else {
		return nil, e.ErrInvalidMatch
	}
}

func (d *Dgraph) UpdateProfile(ctx context.Context, user *models.User) error {
	user.DType = []string{"User"}
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	mutation := &api.Mutation{SetJson: json, CommitNow: true}
	_, err = d.client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dgraph) DeleteProfileAttribute(ctx context.Context, uid, attribute, value string) error {
	if !slices.Contains(attributes, attribute) {
		return e.ErrUnknownPredicate
	}
	nquads := []byte(fmt.Sprintf("<%s> <%s> %q .", uid, attribute, value))
	mutation := &api.Mutation{DelNquads: nquads, CommitNow: true}
	_, err := d.client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return err
	}
	return nil
}
