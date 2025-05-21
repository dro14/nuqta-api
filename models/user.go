package models

type User struct {
	Version     string   `json:"version,omitempty"`
	Uid         string   `json:"uid,omitempty"`
	FirebaseUid string   `json:"firebase_uid,omitempty"`
	Email       string   `json:"email,omitempty"`
	Registered  int64    `json:"registered,omitempty"`
	InvitedBy   *User    `json:"invited_by,omitempty"`
	Name        string   `json:"name,omitempty"`
	Username    string   `json:"username,omitempty"`
	Location    string   `json:"location,omitempty"`
	Birthday    int64    `json:"birthday,omitempty"`
	Color       string   `json:"color,omitempty"`
	Bio         string   `json:"bio,omitempty"`
	Banner      string   `json:"banner,omitempty"`
	Avatars     []string `json:"avatars,omitempty"`
	Thumbnails  []string `json:"thumbnails,omitempty"`
	Followers   int      `json:"followers,omitempty"`
	Invites     int      `json:"invites,omitempty"`
	Blockers    int      `json:"blockers,omitempty"`
	Following   int      `json:"following,omitempty"`
	Posts       int      `json:"posts,omitempty"`
	Replies     int      `json:"replies,omitempty"`
	Media       int      `json:"media,omitempty"`
	Reposts     int      `json:"reposts,omitempty"`
	Likes       int      `json:"likes,omitempty"`
	IsFollowing bool     `json:"is_following,omitempty"`
	IsFollower  bool     `json:"is_follower,omitempty"`
	IsBlocking  bool     `json:"is_blocking,omitempty"`
	IsBlocker   bool     `json:"is_blocker,omitempty"`
	ChatUid     string   `json:"chat_uid,omitempty"`
	Score       int      `json:"score,omitempty"`
}
