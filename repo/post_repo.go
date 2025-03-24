package repo

import (
	"blog-api/db"
	"blog-api/types"
)

func CreatePost(post types.Post) (types.Post, error) {
	conn := db.DB
	var id int

	err := conn.QueryRow(
		"INSERT INTO posts (title, content) VALUES ($1, $2) RETURNING id", 
		post.Title, post.Content).Scan(&id)
	if err != nil { return types.Post{}, err }

	post.ID = id
	return post, nil
}