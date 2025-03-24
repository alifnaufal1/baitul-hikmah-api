ALTER TABLE "posts"
ADD COLUMN category_id INT REFERENCES "categories" (id); 