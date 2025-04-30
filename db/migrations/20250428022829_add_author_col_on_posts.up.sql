ALTER TABLE posts 
ADD COLUMN author_id INT REFERENCES "users" (id);