-- name: CreateFeedFollow :one
With inserted_feed_follow AS (
    INSERT INTO feed_follows(user_id,feed_id) VALUES($1,$2) RETURNING *
)

SELECT inserted_feed_follow.*, feeds.name AS feed_name, users.name AS user_name FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollow :many

SELECT f.name,u.name FROM feed_follows AS ff
 INNER JOIN feeds AS f ON ff.feed_id = f.id
INNER JOIN users AS u ON f.user_id = u.id
WHERE ff.user_id = $1;


-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;
