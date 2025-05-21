CREATE TABLE "blacklist_tokens" (
  token TEXT PRIMARY KEY,
  user_id INT REFERENCES "users" (id)
);