package types

type Post struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Category  Category `json:"category"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

