package models

type User struct {
	DType       []string `json:"dgraph.type,omitempty"`
	Uid         string   `json:"uid,omitempty"`
	FirebaseUid string   `json:"firebase_uid,omitempty"`
	Email       string   `json:"email,omitempty"`
	Name        string   `json:"name,omitempty"`
	Username    string   `json:"username,omitempty"`
	Bio         string   `json:"bio,omitempty"`
	JoinedAt    int64    `json:"joined_at,omitempty"`
	Birthday    int64    `json:"birthday,omitempty"`
	Banner      string   `json:"banner,omitempty"`
	Avatars     []string `json:"avatars,omitempty"`
	Thumbnails  []string `json:"thumbnails,omitempty"`
	Posts       int      `json:"posts,omitempty"`
	Following   int      `json:"following,omitempty"`
	Followers   int      `json:"followers,omitempty"`
	IsFollowing bool     `json:"is_following,omitempty"`
	IsFollowed  bool     `json:"is_followed,omitempty"`
}

type Post struct {
	DType        []string `json:"dgraph.type,omitempty"`
	Uid          string   `json:"uid,omitempty"`
	Text         string   `json:"text,omitempty"`
	PostedAt     int64    `json:"posted_at,omitempty"`
	Author       *User    `json:"author,omitempty"`
	ReplyControl string   `json:"reply_control,omitempty"`
	InReplyTo    *Post    `json:"in_reply_to,omitempty"`
	RepostedBy   *User    `json:"reposted_by,omitempty"`
	Reposted     []*User  `json:"reposted,omitempty"`
	Replies      int      `json:"replies,omitempty"`
	Reposts      int      `json:"reposts,omitempty"`
	Likes        int      `json:"likes,omitempty"`
	Clicks       int      `json:"clicks,omitempty"`
	Views        int      `json:"views,omitempty"`
	Saves        int      `json:"saves,omitempty"`
	Removes      int      `json:"removes,omitempty"`
	IsReplied    bool     `json:"is_replied,omitempty"`
	IsReposted   bool     `json:"is_reposted,omitempty"`
	IsLiked      bool     `json:"is_liked,omitempty"`
	IsClicked    bool     `json:"is_clicked,omitempty"`
	IsViewed     bool     `json:"is_viewed,omitempty"`
	IsSaved      bool     `json:"is_saved,omitempty"`
	Score        float64  `json:"score,omitempty"`
}

type Request struct {
	Uid       string `json:"uid"`
	UserUid   string `json:"user_uid"`
	PostUid   string `json:"post_uid"`
	Username  string `json:"username"`
	Query     string `json:"query"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
	Tab       string `json:"tab"`
	Before    int64  `json:"before"`
	Offset    int    `json:"offset"`
	Edge      string `json:"edge"`
	Target    string `json:"target"`
}
