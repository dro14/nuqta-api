package dgraph

import (
	"context"

	"github.com/dgraph-io/dgo/v240/protos/api"
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

const schema = `
firebase_uid: string @index(hash) .
email: string .
name: string .
username: string .
bio: string .
joined_at: int .
birthday: int .
banner: string .
avatar_th: string .
avatars: [string] .
follow: [uid] @count @reverse .
like: [uid] @count @reverse .
repost: [uid] @count @reverse .
click: [uid] @count @reverse .

text: string .
posted_at: int @index(int) .
author: uid @count @reverse .
in_reply_to: uid @count @reverse .
viewed_by: [uid] @count .

type User {
	firebase_uid
	email
	name
	username
	bio
	joined_at
	birthday
	banner
	avatar_th
	avatars
	follow
	like
	repost
	click
}

type Post {
	text
	posted_at
	author
	in_reply_to
	viewed_by
}`
