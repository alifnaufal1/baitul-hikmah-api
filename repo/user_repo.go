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
    INSERT INTO users (username, password, role)
    VALUES ($1, $2, $3)
    RETURNING id, created_at`, username, hashPassword, "user").Scan(&id, &createdAt)
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
    SELECT id, username, password, role
    FROM users
    WHERE username = $1`, 
    username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
    if err != nil {return types.User{}, err}

    return user, nil
}

func GetUserById(id int) (types.User, error) {
    conn := db.DB

    var user types.User

    err := conn.QueryRow(`
    SELECT id, username, role, url_profile_img
    FROM users
    WHERE id = $1`, 
    id).Scan(&user.ID, &user.Username, &user.Role, &user.URLProfileImg)
    if err != nil {return types.User{}, err}

    return user, nil
}

func GetUsernameById(id int) (types.User, error) {
    conn := db.DB

    var user types.User

    err := conn.QueryRow(`
    SELECT id, username
    FROM users
    WHERE id = $1`, 
    id).Scan(&user.ID, &user.Username)
    if err != nil {return types.User{}, err}

    return user, nil
}

func AddProfileImage(id int, profileImage string) (string, error) {
    conn := db.DB

    err := conn.QueryRow(`
    UPDATE users
    SET url_profile_img = $1
    WHERE id = $2
    RETURNING id`, profileImage, id).Scan(&id)
    if err != nil {return "", err}

    return profileImage, nil
}