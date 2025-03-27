package repo

import (
	"blog-api/db"
	"blog-api/types"
	"blog-api/utils"
)

func CreateUser(username string, hashPassword string) (types.RegisterResponse, error) {
    conn := db.DB
    var id int
    var createdAt string

    err := conn.QueryRow(`
    INSERT INTO users (username, password, is_admin)
    VALUES ($1, $2, $3)
    RETURNING id, created_at`, username, hashPassword, false).Scan(&id, &createdAt)
    if err != nil {return types.RegisterResponse{}, err}

    registerResponse := types.RegisterResponse{
        ID: id,
        Username: username,
        RegisterAT: utils.FromTimestamp(createdAt),
    }

    return registerResponse, nil
}

func GetUserByUsername(username string) (types.User, error) {
    conn := db.DB

    var user types.User
    err := conn.QueryRow(`
    SELECT id, username, password, is_admin
    FROM users
    WHERE username = $1`, 
    username).Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin)
    if err != nil {return types.User{}, err}

    return user, nil
}