package types

type Post struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Description string   `json:"description"`
	CategoryID  int      `json:"category_id"`
	AuthorID    int      `json:"author_id"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type PostResponse struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Description string   `json:"description"`
	PostImg     string   `json:"post_img"`
	Category    string   `json:"category"`
	Author      string   `json:"author"`
	AuthorImg   string   `json:"author_img"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"created_at"`
}