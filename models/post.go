package models

type Post struct {
	Uid              string   `json:"uid,omitempty"`
	Timestamp        int64    `json:"timestamp,omitempty"`
	WhoCanReply      string   `json:"who_can_reply,omitempty"`
	Text             string   `json:"text,omitempty"`
	Images           []string `json:"images,omitempty"`
	Edited           int64    `json:"edited,omitempty"`
	Hidden           int64    `json:"hidden,omitempty"`
	Author           *User    `json:"author,omitempty"`
	InReplyTo        *Post    `json:"in_reply_to,omitempty"`
	RepostedBy       *User    `json:"reposted_by,omitempty"`
	Replies          int      `json:"replies,omitempty"`
	Reposts          int      `json:"reposts,omitempty"`
	Likes            int      `json:"likes,omitempty"`
	Clicks           int      `json:"clicks,omitempty"`
	Views            int      `json:"views,omitempty"`
	Saves            int      `json:"saves,omitempty"`
	IsReplied        bool     `json:"is_replied,omitempty"`
	IsReposted       bool     `json:"is_reposted,omitempty"`
	IsLiked          bool     `json:"is_liked,omitempty"`
	IsClicked        bool     `json:"is_clicked,omitempty"`
	IsViewed         bool     `json:"is_viewed,omitempty"`
	IsSaved          bool     `json:"is_saved,omitempty"`
	Score            int      `json:"score,omitempty"`
	RepostsTimestamp int64    `json:"reposts|timestamp,omitempty"`
}
