-- name: CreateFeed :one
INSERT INTO feeds (name,url,user_id) 
VALUES ($1,$2,$3)
RETURNING * ;


-- name: GetFeeds :many
SELECT name,url,user_id FROM feeds;

-- name: GetFeedByUrl :one
SELECT id FROM feeds WHERE url = $1;