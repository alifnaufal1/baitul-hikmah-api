package repo

import (
	"blog-api/db"
	"blog-api/types"
)

func CreateCategory(category types.Category) (types.Category, error) {
  conn := db.DB
	var id int

	err := conn.QueryRow(`
  INSERT INTO categories (name, description) 
  VALUES ($1, $2) 
  RETURNING id`,
	category.Name, category.Description).Scan(&id)
	
  if err != nil {return types.Category{}, err}

  category.ID = id
  return category, nil
}

func GetAllCategory() ([]types.Category, error) {
  conn :=db.DB

  rows, err := conn.Query(`
  SELECT id, name, description 
  FROM categories 
  WHERE deleted_at IS NULL
  ORDER BY created_at DESC`)
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

func GetCategoryById(id int) (types.Category, error) {
  conn := db.DB
  var category types.Category

  err := conn.QueryRow(`
  SELECT id, name, description 
  FROM categories 
  WHERE id = $1 AND deleted_at IS NULL`, 
  id).Scan(&category.ID, &category.Name, &category.Description)
  if err != nil {return types.Category{}, err}

  return category, nil
}

func UpdateCategory(id int, category types.Category) (types.Category, error) {
  conn := db.DB

  err := conn.QueryRow(`
  UPDATE categories 
  SET 
    name = $1, 
    description = $2 
  WHERE id = $3 AND deleted_at IS NULL 
  RETURNING id`,
  category.Name, category.Description, id).Scan(&id)
  if err != nil {return types.Category{}, err}

  category.ID = id
  return category, nil
}

func DeleteCategory(id int) (error) {
  conn := db.DB

  err := conn.QueryRow(`
  UPDATE categories
  SET deleted_at = NOW()
  WHERE id = $1 AND deleted_at IS NULL
  RETURNING id`, id).Scan(&id)
  if err != nil {return err}

  return nil
}