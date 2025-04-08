package data

import (
	"context"

	"github.com/dgraph-io/dgo/v240/protos/api"
)

func (d *Data) GetSchema(ctx context.Context) (string, error) {
	query := "schema {}"
	bytes, err := d.graphGet(ctx, query, nil)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (d *Data) UpdateSchema(ctx context.Context) error {
	operation := &api.Operation{Schema: GraphSchema}
	err := d.graph.Alter(ctx, operation)
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) DeleteSchema(ctx context.Context) error {
	operation := &api.Operation{DropAll: true}
	err := d.graph.Alter(ctx, operation)
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) DeletePredicate(ctx context.Context, predicate string) error {
	operation := &api.Operation{DropOp: api.Operation_ATTR, DropValue: predicate}
	err := d.graph.Alter(ctx, operation)
	if err != nil {
		return err
	}
	return nil
}

const GraphSchema = `
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

const DdSchema = `
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    firebase_uid VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    registered BIGINT NOT NULL,
    name VARCHAR(255),
    username VARCHAR(255) UNIQUE NOT NULL,
    location VARCHAR(255),
    birthday BIGINT,
    color VARCHAR(50),
    bio TEXT,
    banner VARCHAR(255),
    avatars VARCHAR(255)[],
    thumbnails VARCHAR(255)[],
    search_vector TSVECTOR NOT NULL
);

CREATE INDEX users_firebase_uid_idx ON users(firebase_uid);
CREATE INDEX users_username_lower_idx ON users(LOWER(username));
CREATE INDEX users_search_idx ON users USING GIN(search_vector);

CREATE TABLE posts (
    id VARCHAR(255) PRIMARY KEY,
    timestamp BIGINT NOT NULL,
    text TEXT NOT NULL,
    who_can_reply VARCHAR(50) NOT NULL 
);`
