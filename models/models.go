package models

type User struct {
	DType           []string `json:"dgraph.type,omitempty"`
	Uid             string   `json:"uid,omitempty"`
	Name            string   `json:"name,omitempty"`
	Username        string   `json:"username,omitempty"`
	Bio             string   `json:"bio,omitempty"`
	Birthday        int      `json:"birthday,omitempty"`
	JoinedAt        int      `json:"joined_at,omitempty"`
	Banner          string   `json:"banner,omitempty"`
	Avatars         []string `json:"avatars,omitempty"`
	Posts           int      `json:"posts,omitempty"`
	Following       int      `json:"following,omitempty"`
	Followers       int      `json:"followers,omitempty"`
	Email           string   `json:"email,omitempty"`
	IsEmailVerified bool     `json:"is_email_verified,omitempty"`
	IsAnonymous     bool     `json:"is_anonymous,omitempty"`
	PhoneNumber     string   `json:"phone_number,omitempty"`
	ProviderId      string   `json:"provider_id,omitempty"`
	ProviderUid     string   `json:"provider_uid,omitempty"`
	FirebaseUid     string   `json:"firebase_uid,omitempty"`
}

type Post struct {
	DType        []string `json:"dgraph.type,omitempty"`
	Uid          string   `json:"uid,omitempty"`
	Text         string   `json:"text,omitempty"`
	PostedAt     int      `json:"posted_at,omitempty"`
	AuthorUid    string   `json:"author_uid,omitempty"`
	InReplyToUid string   `json:"in_reply_to_uid,omitempty"`
	Author       *User    `json:"author,omitempty"`
	InReplyTo    *Post    `json:"in_reply_to,omitempty"`
	Views        int      `json:"views,omitempty"`
	Likes        int      `json:"likes,omitempty"`
	Reposts      int      `json:"reposts,omitempty"`
	Replies      int      `json:"replies,omitempty"`
	Clicks       int      `json:"clicks,omitempty"`
}
