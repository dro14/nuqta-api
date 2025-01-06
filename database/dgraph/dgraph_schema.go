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
name: string .
username: string @index(hash) .
bio: string .
birthday: int .
joined_at: int .
banner: string .
avatars: [string] .
email: string .
is_email_verified: bool .
is_anonymous: bool .
phone_number: string .
provider_id: string .
provider_uid: string .
firebase_uid: string @index(hash) .
follow: [uid] @count @reverse .
like: [uid] @count @reverse .
repost: [uid] @count @reverse .
click: [uid] @count @reverse .

text: string .
posted_at: int @index(int) .
author_uid: uid @count @reverse .
in_reply_to_uid: uid @count @reverse .
viewed_by: [uid] @count .

type User {
	name
	username
	bio
	birthday
	joined_at
	banner
	avatars
	email
	is_email_verified
	is_anonymous
	phone_number
	provider_id
	provider_uid
	firebase_uid
	follow
	like
	repost
	click
}

type Post {
	text
	posted_at
	author_uid
	in_reply_to_uid
	viewed_by
}`
