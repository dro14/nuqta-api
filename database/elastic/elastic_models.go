package elastic

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	HitCount int    `json:"hit_count"`
}
