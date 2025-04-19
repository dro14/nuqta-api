package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dro14/nuqta-service/models"
	"github.com/lib/pq"
)

func (d *Data) CreateChat(ctx context.Context, uid, chatWith string) (*models.Chat, error) {
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
		return nil, err
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
		return nil, err
	}

	chat := &models.Chat{
		Uid:      chatUid,
		ChatWith: chatWith,
	}
	return chat, nil
}

func (d *Data) GetChats(ctx context.Context, uid, type_ string) ([]*models.Chat, error) {
	vars := map[string]string{
		"$uid":  uid,
		"$type": type_,
	}
	bytes, err := d.graphGet(ctx, chatsQuery, vars)
	if err != nil {
		return nil, err
	}
	var response map[string][]map[string][]*models.Chat
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}
	chats := make([]*models.Chat, 0)
	for _, user := range response["users"] {
		for _, chat := range user["chats"] {
			chatWith := "yordamchi"
			if len(chat.Members) == 1 {
				chatWith = chat.Members[0].Uid
			}
			chats = append(chats, &models.Chat{
				Uid:      chat.Uid,
				ChatWith: chatWith,
			})
		}
	}
	return chats, nil
}

func (d *Data) GetUpdates(ctx context.Context, type_ string, chatUids []string, after int64) ([]*models.Message, error) {
	query := "SELECT id, timestamp, chat_uid, author_uid, in_reply_to, text, images FROM yordamchi_messages WHERE chat_uid = ANY($1) AND timestamp > $2 ORDER BY timestamp"
	if type_ == "private" {
		query = "SELECT id, timestamp, chat_uid, author_uid, in_reply_to, text, images, viewed, edited, deleted FROM private_messages WHERE chat_uid = ANY($1) AND timestamp > $2 ORDER BY timestamp"
	}
	rows, err := d.dbQuery(ctx, query, pq.Array(chatUids), after)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]*models.Message, 0)
	for rows.Next() {
		message := &models.Message{}
		var nullInReplyTo sql.NullInt64
		var nullText sql.NullString
		var nullViewed sql.NullInt64
		var nullEdited sql.NullInt64
		var nullDeleted sql.NullInt64
		dest := []any{&message.Id, &message.Timestamp, &message.ChatUid, &message.AuthorUid, &nullInReplyTo, &nullText, pq.Array(&message.Images)}
		if type_ == "private" {
			dest = append(dest, &nullViewed, &nullEdited, &nullDeleted)
		}
		err = rows.Scan(dest...)
		if err != nil {
			log.Printf("can't scan message after %d: %s", after, err)
			continue
		}
		if nullInReplyTo.Valid {
			message.InReplyTo = nullInReplyTo.Int64
		}
		if nullText.Valid {
			message.Text = nullText.String
		}
		if nullViewed.Valid {
			message.Viewed = nullViewed.Int64
		}
		if nullEdited.Valid {
			message.Edited = nullEdited.Int64
		}
		if nullDeleted.Valid {
			message.Deleted = nullDeleted.Int64
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (d *Data) CreateMessage(ctx context.Context, message *models.Message, type_ string) error {
	message.Timestamp = time.Now().Unix()
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
		fmt.Sprintf("INSERT INTO %s_messages (timestamp, chat_uid, author_uid, in_reply_to, text, images) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", type_),
		[]any{message.Timestamp, message.ChatUid, message.AuthorUid, nullInReplyTo, nullText, pq.Array(message.Images)},
		&message.Id,
	)
}

func (d *Data) GetMessages(ctx context.Context, type_, chatUid string, before int64) ([]*models.Message, error) {
	query := "SELECT id, timestamp, chat_uid, author_uid, in_reply_to, text, images FROM yordamchi_messages WHERE chat_uid = $1 AND timestamp < $2 ORDER BY timestamp DESC LIMIT 20"
	if type_ == "private" {
		query = "SELECT id, timestamp, chat_uid, author_uid, in_reply_to, text, images, viewed, edited, deleted FROM private_messages WHERE chat_uid = $1 AND timestamp < $2 ORDER BY timestamp DESC LIMIT 20"
	}
	rows, err := d.dbQuery(ctx, query, chatUid, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]*models.Message, 0)
	for rows.Next() {
		message := &models.Message{}
		var nullInReplyTo sql.NullInt64
		var nullText sql.NullString
		var nullViewed sql.NullInt64
		var nullEdited sql.NullInt64
		var nullDeleted sql.NullInt64
		dest := []any{&message.Id, &message.Timestamp, &message.ChatUid, &message.AuthorUid, &nullInReplyTo, &nullText, pq.Array(&message.Images)}
		if type_ == "private" {
			dest = append(dest, &nullViewed, &nullEdited, &nullDeleted)
		}
		err = rows.Scan(dest...)
		if err != nil {
			log.Printf("can't scan message from chat %s before %d: %s", chatUid, before, err)
			continue
		}
		if nullInReplyTo.Valid {
			message.InReplyTo = nullInReplyTo.Int64
		}
		if nullText.Valid {
			message.Text = nullText.String
		}
		if nullViewed.Valid {
			message.Viewed = nullViewed.Int64
		}
		if nullEdited.Valid {
			message.Edited = nullEdited.Int64
		}
		if nullDeleted.Valid {
			message.Deleted = nullDeleted.Int64
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (d *Data) ViewPrivateMessage(ctx context.Context, message *models.Message) error {
	message.Viewed = time.Now().Unix()
	return d.dbExec(ctx,
		"UPDATE private_messages SET viewed = $1 WHERE id = $2 AND author_uid = $3",
		message.Viewed, message.Id, message.AuthorUid,
	)
}

func (d *Data) EditPrivateMessage(ctx context.Context, message *models.Message) error {
	message.Edited = time.Now().Unix()
	return d.dbExec(ctx,
		"UPDATE private_messages SET text = $1, edited = $2 WHERE id = $3 AND author_uid = $4",
		message.Text, message.Edited, message.Id, message.AuthorUid,
	)
}

func (d *Data) DeletePrivateMessage(ctx context.Context, message *models.Message) error {
	message.Deleted = time.Now().Unix()
	return d.dbExec(ctx,
		"UPDATE private_messages SET deleted = $1 WHERE id = $2 AND author_uid = $3",
		message.Deleted, message.Id, message.AuthorUid,
	)
}

func (d *Data) UpdateYordamchiMessage(ctx context.Context, message *models.Message) error {
	message.Timestamp = time.Now().Unix()
	return d.dbExec(ctx,
		"UPDATE yordamchi_messages SET timestamp = $1, author_uid = $2, text = $3 WHERE id = $4",
		message.Timestamp, message.AuthorUid, message.Text, message.Id,
	)
}
