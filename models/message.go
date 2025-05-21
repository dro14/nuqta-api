package models

type Message struct {
	Id           int64    `json:"id,omitempty"`
	Timestamp    int64    `json:"timestamp,omitempty"`
	ChatUid      string   `json:"chat_uid,omitempty"`
	AuthorUid    string   `json:"author_uid,omitempty"`
	RecipientUid string   `json:"recipient_uid,omitempty"`
	InReplyTo    int64    `json:"in_reply_to,omitempty"`
	Text         string   `json:"text,omitempty"`
	Images       []string `json:"images,omitempty"`
	Viewed       int64    `json:"viewed,omitempty"`
	Liked        int64    `json:"liked,omitempty"`
	Edited       int64    `json:"edited,omitempty"`
	Deleted      int64    `json:"deleted,omitempty"`
}
