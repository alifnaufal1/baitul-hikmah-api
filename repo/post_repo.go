package repo

import (
	"blog-api/db"
	"blog-api/types"
	"blog-api/utils"
	"database/sql"
)

func CreatePost(post types.Post) (types.PostResponse, error) {
	conn := db.DB
	var id int
	var createdAt string

	category, err := GetCategoryById(post.CategoryID)
	if err != nil {
		return types.PostResponse{}, err
	}

	err = conn.QueryRow(`
	INSERT INTO posts (title, content, category_id) 
	VALUES ($1, $2, $3) 
	RETURNING id, created_at`, 
	post.Title, post.Content, post.CategoryID).Scan(&id, &createdAt)
	if err != nil {return types.PostResponse{}, err}

	postResponse := types.PostResponse{
		ID:        id,
		Title:     post.Title,
		Content:   post.Content,
		Category:  category.Name,
		CreatedAt: utils.FromTimestamp(createdAt),
	}

	return postResponse, nil
}

func GetAllPost(categoryId int) ([]types.PostResponse, error) {
	conn := db.DB
	var rows *sql.Rows
	var err error

	if categoryId != 0 {
		rows, err = conn.Query(`
		SELECT id, title, content, category_id, created_at
		FROM posts 
		WHERE deleted_at IS NULL AND category_id = $1`, categoryId)
	} else {
		rows, err = conn.Query(`
		SELECT id, title, content, category_id, created_at
		FROM posts 
		WHERE deleted_at IS NULL`)
	}

	if err != nil {return []types.PostResponse{}, err}
	defer rows.Close()

	var posts []types.PostResponse
	for rows.Next() {
		post := types.PostResponse{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
		if err != nil {return []types.PostResponse{}, err}
		post.CreatedAt = utils.FromTimestamp(post.CreatedAt)
		posts = append(posts, post)  
	}
	
	if err = rows.Err(); err != nil {return []types.PostResponse{}, err}
	if len(posts) == 0 {return nil, sql.ErrNoRows}
	return posts, nil
}

func GetPostById(id int) (types.PostResponse, error) {
	conn := db.DB

	var post types.Post
	err := conn.QueryRow(`
	SELECT id, title, content, category_id, created_at 
	FROM posts 
	WHERE id = $1 AND deleted_at IS NULL`,
	id).Scan(&post.ID, &post.Title, &post.Content, &post.CategoryID, &post.CreatedAt)
	if err != nil {return types.PostResponse{}, err}
	
	category, err := GetCategoryById(post.CategoryID)
	if err != nil {return types.PostResponse{}, err}

	postResponse := types.PostResponse{
		ID: id,
		Title: post.Title,
		Content: post.Content,
		Category: category.Name,
		CreatedAt: utils.FromTimestamp(post.CreatedAt),
	}

	return postResponse, nil
}

func UpdatePost(id int, post types.Post) (types.PostResponse, error) {
	conn := db.DB
	
	category, err := GetCategoryById(post.CategoryID)
	if err != nil {return types.PostResponse{}, utils.ErrCategoryNotFound}

	err = conn.QueryRow(`
	UPDATE posts
	SET
	title = $1,
	content = $2,
	category_id = $3
	WHERE id = $4 AND deleted_at IS NULL
	RETURNING id `, 
	post.Title, post.Content, post.CategoryID, id).Scan(&id)
	if err != nil {return types.PostResponse{}, err}

	postResponse := types.PostResponse{
		ID: id,
		Title: post.Title,
		Content: post.Content,
		Category: category.Name,
		CreatedAt: utils.FromTimestamp(post.CreatedAt),
	}
	return postResponse, nil
}

func DeletePost(id int) error {
	conn := db.DB

	err := conn.QueryRow(`
	UPDATE posts 
	SET deleted_at = NOW() 
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING id`, id).Scan(&id)
	if err != nil {return err}

	return nil
}