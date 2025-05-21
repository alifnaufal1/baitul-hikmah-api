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

type UserUpdateRequest struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdateResponse struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	URLProfileImg string `json:"url_profile_img"`
	UpdatedAt     string `json:"updated_at"`
}
