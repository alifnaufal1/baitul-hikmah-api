package types

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type RegisterResponse struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	RegisterAT string `json:"register_at"`
}