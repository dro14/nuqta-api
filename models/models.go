package models

type User struct {
	Name      string   `json:"name,omitempty" bson:"name,omitempty"`
	Username  string   `json:"username,omitempty" bson:"username,omitempty"`
	Bio       string   `json:"bio,omitempty" bson:"bio,omitempty"`
	Birthday  int64    `json:"birthday,omitempty" bson:"birthday,omitempty"`
	CreatedAt int64    `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Banner    string   `json:"banner,omitempty" bson:"banner,omitempty"`
	Avatars   []string `json:"avatars,omitempty" bson:"avatars,omitempty"`

	ID              string `json:"_id,omitempty" bson:"_id,omitempty"`
	Email           string `json:"email,omitempty" bson:"email,omitempty"`
	IsEmailVerified bool   `json:"is_email_verified" bson:"is_email_verified"`
	IsAnonymous     bool   `json:"is_anonymous" bson:"is_anonymous"`
	PhoneNumber     string `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	ProviderID      string `json:"provider_id,omitempty" bson:"provider_id,omitempty"`
	ProviderUID     string `json:"provider_uid,omitempty" bson:"provider_uid,omitempty"`
}
