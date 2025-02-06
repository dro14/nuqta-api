package models

type User struct {
	DType       []string `json:"dgraph.type,omitempty"`
	Uid         string   `json:"uid,omitempty"`
	FirebaseUid string   `json:"firebase_uid,omitempty"`
	Email       string   `json:"email,omitempty"`
	Name        string   `json:"name,omitempty"`
	Username    string   `json:"username,omitempty"`
	Bio         string   `json:"bio,omitempty"`
	JoinedAt    int      `json:"joined_at,omitempty"`
	Birthday    int      `json:"birthday,omitempty"`
	Banner      string   `json:"banner,omitempty"`
	Avatars     []string `json:"avatars,omitempty"`
	Thumbnails  []string `json:"thumbnails,omitempty"`
	Posts       int      `json:"posts,omitempty"`
	Following   int      `json:"following,omitempty"`
	Followers   int      `json:"followers,omitempty"`
	IsFollowed  bool     `json:"is_followed,omitempty"`
}

type Post struct {
	DType      []string `json:"dgraph.type,omitempty"`
	Uid        string   `json:"uid,omitempty"`
	Text       string   `json:"text,omitempty"`
	PostedAt   int      `json:"posted_at,omitempty"`
	Author     *User    `json:"author,omitempty"`
	InReplyTo  *Post    `json:"in_reply_to,omitempty"`
	Views      int      `json:"views,omitempty"`
	Likes      int      `json:"likes,omitempty"`
	Reposts    int      `json:"reposts,omitempty"`
	Replies    int      `json:"replies,omitempty"`
	Clicks     int      `json:"clicks,omitempty"`
	IsLiked    bool     `json:"is_liked,omitempty"`
	IsReposted bool     `json:"is_reposted,omitempty"`
}
