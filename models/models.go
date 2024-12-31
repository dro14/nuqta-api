package models

type User struct {
	Name     string   `json:"name,omitempty"`
	Username string   `json:"username,omitempty"`
	Bio      string   `json:"bio,omitempty"`
	Birthday int64    `json:"birthday,omitempty"`
	JoinedAt int64    `json:"joined_at,omitempty"`
	Banner   string   `json:"banner,omitempty"`
	Avatars  []string `json:"avatars,omitempty"`

	Email           string `json:"email,omitempty"`
	IsEmailVerified bool   `json:"is_email_verified"`
	IsAnonymous     bool   `json:"is_anonymous"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	ProviderID      string `json:"provider_id,omitempty"`
	ProviderUID     string `json:"provider_uid,omitempty"`
	FirebaseUID     string `json:"firebase_uid,omitempty"`
	UID             string `json:"uid,omitempty"`
}

type Post struct {
	CreatedAt int64  `json:"created_at,omitempty"`
	Text      string `json:"text,omitempty"`
}
