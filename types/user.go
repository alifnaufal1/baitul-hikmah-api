package types

type User struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Role          string `json:"role"`
	URLProfileImg string `json:"url_profile_img"`
}

type RegisterResponse struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	RegisterAT string `json:"register_at"`
}