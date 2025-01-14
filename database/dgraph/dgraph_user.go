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

var predicates = []string{
	"name",
	"bio",
	"birthday",
	"banner",
	"avatars",
}

func (d *Dgraph) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
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

func (d *Dgraph) GetUser(ctx context.Context, by, value string) (*models.User, error) {
	function, ok := functions[by]
	if !ok {
		return nil, e.ErrUnknownParam
	}

	query := fmt.Sprintf(userQuery, fmt.Sprintf(function, value))
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.User
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	if len(response["users"]) > 0 && response["users"][0].JoinedAt != 0 {
		return response["users"][0], nil
	} else {
		return nil, e.ErrNotFound
	}
}

func (d *Dgraph) UpdateUser(ctx context.Context, user *models.User) error {
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

func (d *Dgraph) DeleteUserPredicate(ctx context.Context, uid, predicate, value string) error {
	if !slices.Contains(predicates, predicate) {
		return e.ErrUnknownPredicate
	}
	nquads := []byte(fmt.Sprintf("<%s> <%s> <%s> .", uid, predicate, value))
	mutation := &api.Mutation{DelNquads: nquads, CommitNow: true}
	_, err := d.client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dgraph) DeleteUser(ctx context.Context, uid string) error {
	nquads := []byte(fmt.Sprintf("<%s> * * .", uid))
	mutation := &api.Mutation{DelNquads: nquads, CommitNow: true}
	_, err := d.client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return err
	}
	return nil
}
