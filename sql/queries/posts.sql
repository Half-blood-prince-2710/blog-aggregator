-- name: CreatePost :one
INSERT INTO posts(title,url,description,published_at,feed_id)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;