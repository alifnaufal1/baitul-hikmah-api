package repo

import (
	"blog-api/db"
	"blog-api/types"
	"blog-api/utils"
	"database/sql"
	"strconv"
)

func CreatePost(post types.Post, author int) (types.PostResponse, error) {
	conn := db.DB
	var id int
	var createdAt string

	category, err := GetCategoryById(post.CategoryID)
	if err != nil {
		return types.PostResponse{}, err
	}
	
	user, err := GetUserById(author)
	if err != nil {
		return types.PostResponse{}, err
	}

	err = conn.QueryRow(`
	INSERT INTO posts (title, content, description, category_id, author_id) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, created_at`, 
	post.Title, post.Content, post.Description, post.CategoryID, user.ID).Scan(&id, &createdAt)
	if err != nil {return types.PostResponse{}, err}

	postResponse := types.PostResponse{
		ID:        id,
		Title:     post.Title,
		Content:   post.Content,
		Description: post.Description,
		Category:  category.Name,
		Author:    user.Username,
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
		SELECT id, title, content, description, category_id, author_id, created_at
		FROM posts 
		WHERE deleted_at IS NULL AND category_id = $1`, categoryId)
	} else {
		rows, err = conn.Query(`
		SELECT id, title, content, url_post_img, description, category_id, author_id, created_at
		FROM posts 
		WHERE deleted_at IS NULL`)
	}

	if err != nil {return []types.PostResponse{}, err}
	defer rows.Close()

	var posts []types.PostResponse
	for rows.Next() {
		post := types.PostResponse{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.PostImg, &post.Description, &post.Category, &post.Author, &post.CreatedAt)
		if err != nil {return []types.PostResponse{}, err}
		authorId, _ := strconv.Atoi(post.Author)
		user, _ := GetUserById(authorId)
		post.Author = user.Username
		post.AuthorImg = user.URLProfileImg
		timeDiff := utils.GetHumanReadableTimeDiff(post.CreatedAt)
		if timeDiff == "invalid date" {return []types.PostResponse{}, err}
		post.CreatedAt = timeDiff
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
	SELECT id, title, content, category_id, author_id, created_at 
	FROM posts 
	WHERE id = $1 AND deleted_at IS NULL`,
	id).Scan(&post.ID, &post.Title, &post.Content, &post.CategoryID, &post.AuthorID, &post.CreatedAt)
	if err != nil {return types.PostResponse{}, err}
	
	category, err := GetCategoryById(post.CategoryID)
	if err != nil {return types.PostResponse{}, err}

	user, err := GetUserById(post.AuthorID)
	if err != nil {return types.PostResponse{}, err}

	postResponse := types.PostResponse{
		ID: id,
		Title: post.Title,
		Content: post.Content,
		Category: category.Name,
		Author: user.Username,
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

func UpdatePostImage(id int, URLPostImg string) error {
	conn := db.DB

	err := conn.QueryRow(`
	UPDATE posts
	SET url_post_img = $1
	WHERE id = $2 AND deleted_at IS NULL
	RETURNING id`, 
	URLPostImg, id).Scan(&id)
	if err != nil {return err}

	return nil
}