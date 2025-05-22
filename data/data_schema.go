package data

const GraphSchema = `
<author>: uid @count @reverse .
<block>: [uid] @count @reverse .
<chat>: [uid] .
<click>: [uid] @count @reverse .
<follow>: [uid] @count @reverse .
<has_media>: bool .
<hidden>: int .
<in_reply_to>: uid @count @reverse .
<invited_by>: uid @count @reverse .
<like>: [uid] @count @reverse .
<members>: [uid] .
<registered>: int @index(int) .
<report>: [uid] @count @reverse .
<repost>: [uid] @count @reverse .
<save>: [uid] @count @reverse .
<timestamp>: int @index(int) .
<view>: [uid] @count @reverse .

type <post> {
	timestamp
    has_media
	author
	in_reply_to
    hidden
	repost
	like
	click
	view
	save
    report
}

type <private_chat> {
    members
}

type <user> {
	registered
	invited_by
	follow
    chat
    block
}
    
type <yordamchi_chat> {
    members
}`

const DdSchema = `
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    firebase_uid VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    registered BIGINT NOT NULL,
    name VARCHAR(255),
    username VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    birthday BIGINT,
    color VARCHAR(50),
    bio TEXT,
    banner VARCHAR(255),
    avatars VARCHAR(255)[],
    thumbnails VARCHAR(255)[],
    search_vector TSVECTOR NOT NULL
);
CREATE INDEX users_search_idx ON users USING GIN(search_vector);
CREATE UNIQUE INDEX users_username_idx ON users(LOWER(username));

CREATE TABLE posts (
    id VARCHAR(255) PRIMARY KEY,
    timestamp BIGINT NOT NULL,
    text TEXT,
    who_can_reply VARCHAR(50) NOT NULL,
    images VARCHAR(255)[],
    edited BIGINT
);

CREATE TABLE private_messages (
    id BIGSERIAL PRIMARY KEY,
    last_updated BIGINT NOT NULL,
    chat_uid VARCHAR(255) NOT NULL,
    author_uid VARCHAR(255) NOT NULL,
    in_reply_to BIGINT,
    text TEXT,
    images VARCHAR(255)[],
    viewed BIGINT,
    edited BIGINT,
    deleted BIGINT,
    recipient_uid VARCHAR(255) NOT NULL,
    timestamp BIGINT NOT NULL,
    liked BIGINT
);
CREATE INDEX private_messages_chat_uid_idx ON private_messages(chat_uid);
CREATE INDEX private_messages_last_updated_idx ON private_messages(last_updated);
CREATE INDEX private_messages_timestamp_idx ON private_messages(timestamp);
CREATE INDEX private_messages_deleted_idx ON private_messages(deleted);

CREATE TABLE yordamchi_messages (
    id BIGSERIAL PRIMARY KEY,
    timestamp BIGINT NOT NULL,
    chat_uid VARCHAR(255) NOT NULL,
    author_uid VARCHAR(255) NOT NULL,
    in_reply_to BIGINT,
    text TEXT,
    images VARCHAR(255)[],
    deleted BIGINT
);
CREATE INDEX yordamchi_messages_timestamp_idx ON yordamchi_messages(timestamp);
CREATE INDEX yordamchi_messages_chat_uid_idx ON yordamchi_messages(chat_uid);
CREATE INDEX yordamchi_messages_deleted_idx ON yordamchi_messages(deleted);`
