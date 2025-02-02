-- name: CreateFeed :one
INSERT INTO feeds (name,url,user_id) 
VALUES ($1,$2,$3)
RETURNING * ;


-- name: GetFeeds :many
SELECT name,url,user_id FROM feeds;

-- name: GetFeedByUrl :one
SELECT id FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds 
SET updated_at = $1 , last_fetched_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT id, name, url, last_fetched_at
FROM feeds
ORDER BY last_fetched_at NULLS FIRST, created_at ASC
LIMIT 1;