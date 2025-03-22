package dgraph

import (
	"context"

	"github.com/dgraph-io/dgo/v240/protos/api"
)

func (d *Dgraph) GetSchema(ctx context.Context) (string, error) {
	query := "schema {}"
	bytes, err := d.get(ctx, query, nil)
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
registered: int .
name: string .
username: string .
birthday: int .
color: string .
bio: string .
banner: string .
avatars: [string] .
thumbnails: [string] .
follow: [uid] @count @reverse .

text: string .
timestamp: int @index(int) .
who_can_reply: string .
author: uid @count @reverse .
in_reply_to: uid @count @reverse .
repost: [uid] @count @reverse .
like: [uid] @count @reverse .
click: [uid] @count @reverse .
view: [uid] @count @reverse .
save: [uid] @count @reverse .
report: [uid] @count @reverse .

type User {
	firebase_uid
	email
	registered
	name
	username
	birthday
	color
	bio
	banner
	avatars
	thumbnails
	follow
}

type Post {
	text
	timestamp
	who_can_reply
	author
	in_reply_to
	repost
	like
	click
	view
	save
	report
}`
