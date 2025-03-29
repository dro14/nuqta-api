package models

type User struct {
	DType             []string `json:"dgraph.type,omitempty"`
	FirebaseUid       string   `json:"firebase_uid,omitempty"`
	Version           string   `json:"version,omitempty"`
	Uid               string   `json:"uid,omitempty"`
	Email             string   `json:"email,omitempty"`
	Registered        int64    `json:"registered,omitempty"`
	Name              string   `json:"name,omitempty"`
	Username          string   `json:"username,omitempty"`
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
	RepostedTimestamp int64    `json:"reposted|timestamp,omitempty"`
}

type Post struct {
	DType       []string `json:"dgraph.type,omitempty"`
	Uid         string   `json:"uid,omitempty"`
	Text        string   `json:"text,omitempty"`
	Timestamp   int64    `json:"timestamp,omitempty"`
	WhoCanReply string   `json:"who_can_reply,omitempty"`
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

type Request struct {
	Uid       string   `json:"uid"`
	Tab       string   `json:"tab"`
	UserUid   string   `json:"user_uid"`
	UserUids  []string `json:"user_uids"`
	PostUid   string   `json:"post_uid"`
	PostUids  []string `json:"post_uids"`
	Username  string   `json:"username"`
	Query     string   `json:"query"`
	After     string   `json:"after"`
	Before    int64    `json:"before"`
	Offset    int64    `json:"offset"`
	Attribute string   `json:"attribute"`
	Value     string   `json:"value"`
	Source    []string `json:"source"`
	Edge      []string `json:"edge"`
	Target    []string `json:"target"`
}
