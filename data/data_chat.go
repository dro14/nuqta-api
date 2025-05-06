package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/dro14/nuqta-service/models"
	"github.com/lib/pq"
)

func (d *Data) CreateChat(ctx context.Context, uid, chatWith string) (string, error) {
	type_ := "yordamchi_chat"
	members := []map[string]string{{"uid": uid}}
	if chatWith != "yordamchi" {
		type_ = "private_chat"
		members = append(members, map[string]string{"uid": chatWith})
	}
	object := map[string]any{
		"dgraph.type": type_,
		"members":     members,
		"uid":         "_:chat",
	}
	assigned, err := d.graphSet(ctx, object)
	if err != nil {
		return "", err
	}

	chatUid := assigned.Uids["chat"]
	source := []string{uid}
	edge := []string{"chat"}
	target := []string{chatUid}
	if chatWith != "yordamchi" {
		source = append(source, chatWith)
		edge = append(edge, "chat")
		target = append(target, chatUid)
	}
	err = d.CreateEdge(ctx, source, edge, target)
	if err != nil {
		object = map[string]any{"uid": chatUid}
		d.graphDelete(ctx, object)
		return "", err
	}

	return chatUid, nil
}

func (d *Data) GetChats(ctx context.Context, uid, type_ string) ([]string, error) {
	vars := map[string]string{
		"$uid":  uid,
		"$type": type_ + "_chat",
	}
	bytes, err := d.graphGet(ctx, chatsQuery, vars)
	if err != nil {
		return nil, err
	}
	var response map[string][]map[string][]map[string]string
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}
	chatUids := make([]string, 0)
	for _, user := range response["users"] {
		for _, chat := range user["chats"] {
			chatUids = append(chatUids, chat["uid"])
		}
	}
	return chatUids, nil
}

func (d *Data) GetUpdates(ctx context.Context, uid, type_ string, chatUids []string, after int64) ([]*models.Message, error) {
	query := "SELECT id, timestamp, chat_uid, author_uid, in_reply_to, text, images FROM yordamchi_messages WHERE chat_uid = ANY($1) AND timestamp > $2 AND deleted IS NULL"
	if type_ == "private" {
		query = "SELECT id, timestamp, chat_uid, author_uid, in_reply_to, text, images, viewed, edited, deleted, recipient_uid FROM private_messages WHERE chat_uid = ANY($1) AND last_updated > $2 AND (author_uid != $3 OR deleted IS NULL)"
	}
	rows, err := d.dbQuery(ctx, query, pq.Array(chatUids), after, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return decodeMessages(rows, type_), nil
}

func (d *Data) GetMessages(ctx context.Context, uid, type_, chatUid string, before int64) ([]*models.Message, error) {
	query := "SELECT id, timestamp, chat_uid, author_uid, in_reply_to, text, images FROM yordamchi_messages WHERE chat_uid = $1 AND timestamp < $2 AND deleted IS NULL ORDER BY timestamp DESC LIMIT 20"
	if type_ == "private" {
		query = "SELECT id, timestamp, chat_uid, author_uid, in_reply_to, text, images, viewed, edited, deleted, recipient_uid FROM private_messages WHERE chat_uid = $1 AND timestamp < $2 AND (author_uid != $3 OR deleted IS NULL) ORDER BY timestamp DESC LIMIT 20"
	}
	rows, err := d.dbQuery(ctx, query, chatUid, before, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return decodeMessages(rows, type_), nil
}

func (d *Data) CreatePrivate(ctx context.Context, message *models.Message, uid string) error {
	message.Timestamp = time.Now().UnixMilli()
	var nullInReplyTo sql.NullInt64
	var nullText sql.NullString
	if message.InReplyTo != 0 {
		nullInReplyTo.Valid = true
		nullInReplyTo.Int64 = message.InReplyTo
	}
	if message.Text != "" {
		nullText.Valid = true
		nullText.String = message.Text
	}
	return d.dbQueryRow(ctx,
		"INSERT INTO private_messages (timestamp, chat_uid, author_uid, in_reply_to, text, images, recipient_uid, last_updated) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		[]any{message.Timestamp, message.ChatUid, uid, nullInReplyTo, nullText, pq.Array(message.Images), message.RecipientUid, message.Timestamp},
		&message.Id,
	)
}

func (d *Data) ViewPrivate(ctx context.Context, messages []*models.Message, uid string) error {
	now := time.Now().UnixMilli()
	ids := make([]int64, 0)
	for i, message := range messages {
		messages[i].Viewed = now
		ids = append(ids, message.Id)
	}
	return d.dbExec(ctx,
		"UPDATE private_messages SET viewed = $1, last_updated = $2 WHERE id = ANY($3) AND recipient_uid = $4",
		now, now, pq.Array(ids), uid,
	)
}

func (d *Data) EditPrivate(ctx context.Context, message *models.Message, uid string) error {
	message.Edited = time.Now().UnixMilli()
	return d.dbExec(ctx,
		"UPDATE private_messages SET text = $1, images = $2, edited = $3, last_updated = $4 WHERE id = $5 AND author_uid = $6",
		message.Text, pq.Array(message.Images), message.Edited, message.Edited, message.Id, uid,
	)
}

func (d *Data) DeletePrivate(ctx context.Context, message *models.Message, uid string) error {
	message.Deleted = time.Now().UnixMilli()
	return d.dbExec(ctx,
		"UPDATE private_messages SET deleted = $1, last_updated = $2 WHERE id = $3 AND author_uid = $4",
		message.Deleted, message.Deleted, message.Id, uid,
	)
}

func (d *Data) RemovePrivate(ctx context.Context, message *models.Message, uid string) error {
	return d.dbExec(ctx,
		"DELETE FROM private_messages WHERE id = $1 AND recipient_uid = $2",
		message.Id, uid,
	)
}

func (d *Data) CreateYordamchi(ctx context.Context, message *models.Message) error {
	message.Timestamp = time.Now().UnixMilli()
	var nullInReplyTo sql.NullInt64
	var nullText sql.NullString
	if message.InReplyTo != 0 {
		nullInReplyTo.Valid = true
		nullInReplyTo.Int64 = message.InReplyTo
	}
	if message.Text != "" {
		nullText.Valid = true
		nullText.String = message.Text
	}
	return d.dbQueryRow(ctx,
		"INSERT INTO yordamchi_messages (timestamp, chat_uid, author_uid, in_reply_to, text) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		[]any{message.Timestamp, message.ChatUid, message.AuthorUid, nullInReplyTo, nullText},
		&message.Id,
	)
}

func (d *Data) ClearYordamchi(ctx context.Context, message *models.Message) error {
	return d.dbExec(ctx,
		"UPDATE yordamchi_messages SET deleted = $1 WHERE chat_uid = $2 AND id >= $3 AND deleted IS NULL",
		time.Now().UnixMilli(), message.ChatUid, message.Id,
	)
}

func (d *Data) DeleteYordamchi(ctx context.Context, message *models.Message) error {
	message.Deleted = time.Now().UnixMilli()
	return d.dbExec(ctx,
		"UPDATE yordamchi_messages SET deleted = $1 WHERE id = $2 AND author_uid = $3 AND deleted IS NULL",
		message.Deleted, message.Id, message.AuthorUid,
	)
}
