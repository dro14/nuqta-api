package models

type User struct {
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
	Uid             string   `json:"uid,omitempty"`
	DType           []string `json:"dgraph.type,omitempty"`
}

type Post struct {
	CreatedAt int64    `json:"created_at,omitempty"`
	Text      string   `json:"text,omitempty"`
	UID       string   `json:"uid,omitempty"`
	DType     []string `json:"dgraph.type,omitempty"`
}
