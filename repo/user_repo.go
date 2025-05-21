package repo

import (
	"blog-api/db"
	"blog-api/types"
	"blog-api/utils"
	"database/sql"
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
        RegisterAT: utils.GetDate(createdAt),
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

func UpdateUser(user types.UserUpdateRequest) (types.UserUpdateResponse, error) {
    conn := db.DB

    var updatedUser types.UserUpdateResponse

    err := conn.QueryRow(`
    UPDATE users
    SET 
    username = $1, 
    password = $2
    WHERE id = $3
    RETURNING id, updated_at`,
    user.Username, user.Password, user.ID).Scan(&user.ID, &updatedUser.UpdatedAt)
    if err != nil { return types.UserUpdateResponse{}, err }

    updatedUser = types.UserUpdateResponse{
        ID: user.ID,
        Username: user.Username,
        UpdatedAt: utils.GetDateHour(updatedUser.UpdatedAt),
    }
    return updatedUser, nil
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

func AddBlacklistToken(token string, id int) (error) {
    conn := db.DB

    err := conn.QueryRow(`
    INSERT INTO blacklist_tokens (token, user_id)
    VALUES ($1, $2)
    `, token, id).Err()
    if err != nil { return err }

    return nil
}

func IsBlacklistToken(token string, id int) (bool) {
    conn := db.DB

    var blacklistTokens string

    err := conn.QueryRow(`
    SELECT token, user_id
    FROM blacklist_tokens
    WHERE token=$1 AND user_id=$2
    `, token, id).Scan(blacklistTokens)
    
    return err != sql.ErrNoRows

}