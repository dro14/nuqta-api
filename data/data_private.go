package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/dro14/nuqta-service/models"
	"github.com/lib/pq"
)

func (d *Data) CreatePrivateChat(ctx context.Context, uid, userUid string) (string, error) {
	object := map[string]any{
		"dgraph.type": "private_chat",
		"uid":         "_:private_chat",
		"members": []map[string]string{
			{"uid": uid},
			{"uid": userUid},
		},
	}
	assigned, err := d.graphSet(ctx, object)
	if err != nil {
		return "", err
	}
	chatUid := assigned.Uids["private_chat"]
	source := []string{uid, userUid}
	edge := []string{"chat", "chat"}
	target := []string{chatUid, chatUid}
	err = d.CreateEdge(ctx, source, edge, target)
	if err != nil {
		object = map[string]any{"uid": chatUid}
		d.graphDelete(ctx, object)
		return "", err
	}
	return chatUid, nil
}

func (d *Data) GetPrivateChats(ctx context.Context, uid string) ([]*models.Chat, error) {
	vars := map[string]string{
		"$uid": uid,
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
	var chats []*models.Chat
	for _, user := range response["users"] {
		for _, chat := range user["chats"] {
			if len(chat.Members) == 1 {
				chats = append(chats, chat)
			}
		}
	}
	return chats, nil
}

func (d *Data) CreateMessage(ctx context.Context, message *models.Message) error {
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
		"INSERT INTO private_messages (timestamp, chat_uid, author_uid, in_reply_to, text, images) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		[]any{message.Timestamp, message.ChatUid, message.AuthorUid, nullInReplyTo, nullText, pq.Array(message.Images)},
		&message.Id,
	)
}

func (d *Data) GetMessages(ctx context.Context, chatUid string, before int64) ([]*models.Message, error) {
	rows, err := d.dbQuery(ctx,
		"SELECT id, timestamp, chat_uid, author_uid, in_reply_to, text, images FROM private_messages WHERE chat_uid = $1 AND timestamp < $2 ORDER BY timestamp DESC LIMIT 20",
		[]any{chatUid, before},
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		var nullInReplyTo sql.NullInt64
		var nullText sql.NullString
		err = rows.Scan(&message.Id, &message.Timestamp, &message.ChatUid, &message.AuthorUid, &nullInReplyTo, &nullText, pq.Array(&message.Images))
		if err != nil {
			return nil, err
		}
		if nullInReplyTo.Valid {
			message.InReplyTo = nullInReplyTo.Int64
		}
		if nullText.Valid {
			message.Text = nullText.String
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (d *Data) ViewMessage(ctx context.Context, message *models.Message) error {
	message.Viewed = time.Now().Unix()
	return d.dbExec(ctx,
		"UPDATE private_messages SET viewed = $1 WHERE id = $2 AND author_uid != $3",
		[]any{message.Viewed, message.Id, message.AuthorUid},
	)
}

func (d *Data) EditMessage(ctx context.Context, message *models.Message) error {
	message.Edited = time.Now().Unix()
	return d.dbExec(ctx,
		"UPDATE private_messages SET text = $1, edited = $2 WHERE id = $3 AND author_uid = $4",
		[]any{message.Text, message.Edited, message.Id, message.AuthorUid},
	)
}

func (d *Data) DeleteMessage(ctx context.Context, message *models.Message) error {
	message.Deleted = time.Now().Unix()
	return d.dbExec(ctx,
		"UPDATE private_messages SET deleted = $1 WHERE id = $2 AND author_uid = $3",
		[]any{message.Deleted, message.Id, message.AuthorUid},
	)
}
