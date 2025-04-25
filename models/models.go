package models

type User struct {
	Version           string   `json:"version,omitempty"`
	Uid               string   `json:"uid,omitempty"`
	FirebaseUid       string   `json:"firebase_uid,omitempty"`
	Email             string   `json:"email,omitempty"`
	Registered        int64    `json:"registered,omitempty"`
	InvitedBy         *User    `json:"invited_by,omitempty"`
	Name              string   `json:"name,omitempty"`
	Username          string   `json:"username,omitempty"`
	Location          string   `json:"location,omitempty"`
	Birthday          int64    `json:"birthday,omitempty"`
	Color             string   `json:"color,omitempty"`
	Bio               string   `json:"bio,omitempty"`
	Banner            string   `json:"banner,omitempty"`
	Avatars           []string `json:"avatars,omitempty"`
	Thumbnails        []string `json:"thumbnails,omitempty"`
	Followers         int      `json:"followers,omitempty"`
	Following         int      `json:"following,omitempty"`
	Posts             int      `json:"posts,omitempty"`
	Replies           int      `json:"replies,omitempty"`
	Reposts           int      `json:"reposts,omitempty"`
	Likes             int      `json:"likes,omitempty"`
	IsFollowed        bool     `json:"is_followed,omitempty"`
	IsFollowing       bool     `json:"is_following,omitempty"`
	ChatUid           string   `json:"chat_uid,omitempty"`
	RepostedTimestamp int64    `json:"reposted|timestamp,omitempty"`
}

type Post struct {
	Uid         string   `json:"uid,omitempty"`
	Timestamp   int64    `json:"timestamp,omitempty"`
	WhoCanReply string   `json:"who_can_reply,omitempty"`
	Text        string   `json:"text,omitempty"`
	Images      []string `json:"images,omitempty"`
	Author      *User    `json:"author,omitempty"`
	InReplyTo   *Post    `json:"in_reply_to,omitempty"`
	RepostedBy  *User    `json:"reposted_by,omitempty"`
	Reposted    []*User  `json:"reposted,omitempty"`
	Replies     int      `json:"replies,omitempty"`
	Reposts     int      `json:"reposts,omitempty"`
	Likes       int      `json:"likes,omitempty"`
	Clicks      int      `json:"clicks,omitempty"`
	Views       int      `json:"views,omitempty"`
	Saves       int      `json:"saves,omitempty"`
	Reports     int      `json:"reports,omitempty"`
	IsReplied   bool     `json:"is_replied,omitempty"`
	IsReposted  bool     `json:"is_reposted,omitempty"`
	IsLiked     bool     `json:"is_liked,omitempty"`
	IsClicked   bool     `json:"is_clicked,omitempty"`
	IsViewed    bool     `json:"is_viewed,omitempty"`
	IsSaved     bool     `json:"is_saved,omitempty"`
	Score       float64  `json:"score,omitempty"`
}

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
	Edited       int64    `json:"edited,omitempty"`
	Deleted      int64    `json:"deleted,omitempty"`
}
