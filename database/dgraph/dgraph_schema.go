package dgraph

import (
	"context"

	"github.com/dgraph-io/dgo/v240/protos/api"
)

func (d *Dgraph) GetSchema(ctx context.Context) (string, error) {
	query := `schema {}`
	bytes, err := d.get(ctx, query)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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

func (d *Dgraph) DeletePredicate(ctx context.Context, predicate string) error {
	operation := &api.Operation{DropOp: api.Operation_ATTR, DropValue: predicate}
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
avatars: [string] .
thumbnails: [string] .
follow: [uid] @count @reverse .

text: string .
posted_at: int @index(int) .
author: uid @count @reverse .
reply_control: string .
in_reply_to: uid @count @reverse .
repost: [uid] @count @reverse .
like: [uid] @count @reverse .
click: [uid] @count @reverse .
view: [uid] @count @reverse .
remove: [uid] @count @reverse .
save: [uid] @count @reverse .

type User {
	firebase_uid
	email
	name
	username
	bio
	joined_at
	birthday
	banner
	avatars
	thumbnails
	follow
}

type Post {
	text
	posted_at
	author
	reply_control
	in_reply_to
	repost
	like
	click
	view
	remove
	save
}`
