package repo

import (
	"blog-api/db"
	"blog-api/types"
)


func CreateCategory(category types.Category) (types.Category, error) {
  conn := db.DB
	var id int

	err := conn.QueryRow(
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id",
		category.Name, category.Description).Scan(&id)
	
  if err != nil {return types.Category{}, err}

  category.ID = id
  return category, nil
}

func GetAllCategory() ([]types.Category, error) {
  conn :=db.DB

  rows, err := conn.Query(
    "SELECT id, name, description FROM categories ORDER BY created_at DESC")
  if err != nil {return []types.Category{}, err}
  defer rows.Close()

  var categories []types.Category
  for rows.Next() {
    category := types.Category{}
    err := rows.Scan(&category.ID, &category.Name, &category.Description)
    if err != nil {return []types.Category{}, err}
    categories = append(categories, category)  
  }

  if err = rows.Err(); err != nil {return []types.Category{}, err}
  return categories, nil
}

func UpdateCategory(category types.Category) (types.Category, error) {
  conn := db.DB
  var id int

  err := conn.QueryRow(
    `UPDATE categories 
      SET 
        name = $1, 
        description = $2 
      WHERE id = $3 RETURNING id`,
    category.Name, category.Description, category.ID).Scan(&id)

  if err != nil {return types.Category{}, err}

  category.ID = id
  return category, nil
}

func DeleteCategory(id int) (error) {
  conn := db.DB

  _, err := conn.Exec(
    "DELETE FROM categories WHERE id = $1", id)

  if err != nil {return err}

  return nil
}