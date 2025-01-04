package dgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
)

func (d *Dgraph) GetSchema(ctx context.Context) (string, error) {
	query := `schema {}`
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return "", err
	}
	return string(resp.Json), nil
}

func (d *Dgraph) UpdateSchema(ctx context.Context) error {
	operation := &api.Operation{Schema: schema}
	err := d.client.Alter(ctx, operation)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dgraph) DeleteSchema(ctx context.Context) error {
	operation := &api.Operation{DropAll: true}
	err := d.client.Alter(ctx, operation)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dgraph) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.DType = []string{"User"}
	user.Uid = "_:user"
	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	mutation := &api.Mutation{
		SetJson:   bytes,
		CommitNow: true,
	}

	assigned, err := d.client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return nil, err
	}

	user.Uid = assigned.Uids["user"]
	return user, nil
}

func (d *Dgraph) GetUser(ctx context.Context, by, value string) (*models.User, error) {
	var query string
	switch by {
	case "uid":
		query = userByUid
	case "firebase_uid":
		query = userByFirebaseUid
	case "username":
		query = userByUsername
	default:
		return nil, e.ErrUnknownParam
	}

	resp, err := d.client.NewTxn().Query(ctx, fmt.Sprintf(query, value))
	if err != nil {
		return nil, err
	}

	var response map[string][]models.User
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	if len(response["user"]) > 0 {
		return &response["user"][0], nil
	} else {
		return nil, errors.New(by + ":" + value)
	}
}

func (d *Dgraph) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.DType = []string{"User"}
	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	mutation := &api.Mutation{
		SetJson:   bytes,
		CommitNow: true,
	}

	_, err = d.client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *Dgraph) DeleteUser(ctx context.Context, uid string) error {
	mutation := &api.Mutation{
		DeleteJson: []byte(fmt.Sprintf(`{"uid": "%s"}`, uid)),
		CommitNow:  true,
	}

	_, err := d.client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return err
	}

	return nil
}
