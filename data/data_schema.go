package data

const GraphSchema = `
author: uid @count @reverse .
click: [uid] @count @reverse .
follow: [uid] @count @reverse .
in_reply_to: uid @count @reverse .
invited_by: uid @count @reverse .
like: [uid] @count @reverse .
report: [uid] @count @reverse .
repost: [uid] @count @reverse .
save: [uid] @count @reverse .
timestamp: int @index(int) .
view: [uid] @count @reverse .

type User {
	invited_by
	follow
}

type Post {
	timestamp
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
