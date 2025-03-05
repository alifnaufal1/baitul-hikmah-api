How to migrate:
- run on terminal : 
migrate -database "postgres://{db_user}:{db_password}@{db_host}:{port}/blog_db?sslmode=disable" -path db/migrations up