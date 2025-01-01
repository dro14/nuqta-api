package dgraph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/dro14/nuqta-service/models"
)

func (d *Dgraph) ReadSchema(ctx context.Context) (string, error) {
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

func (d *Dgraph) CreateUser(ctx context.Context, user *models.User) error {
	user.DType = []string{"User"}
	user.UID = "_:user"
	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	mutation := &api.Mutation{
		SetJson:   bytes,
		CommitNow: true,
	}

	assigned, err := d.client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return err
	}

	user.UID = assigned.Uids["user"]
	return nil
}

func (d *Dgraph) ReadUser(ctx context.Context, firebaseUid string) (*models.User, error) {
	query := fmt.Sprintf(`
{
	user(func: eq(firebase_uid, "%s")) {
		uid
		expand(_all_)
	}
}`, firebaseUid)

	resp, err := d.client.NewTxn().Query(ctx, query)
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
		return nil, ErrNotFound
	}
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

func (d *Dgraph) DoesUserExist(ctx context.Context, firebaseUid string) (bool, error) {
	query := fmt.Sprintf(`
{
	user(func: eq(firebase_uid, "%s")) {
		uid
	}
}`, firebaseUid)

	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return false, err
	}

	var response map[string][]any
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return false, err
	}

	return len(response["user"]) > 0, nil
}
