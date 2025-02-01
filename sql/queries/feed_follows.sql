-- name: CreateFeedFollow :one
With inserted_feed_follow AS (
    INSERT INTO feed_follows(user_id,feed_id) VALUES($1,$2) RETURNING *
)

SELECT inserted_feed_follow.*, feeds.name AS feed_name, users.name AS user_name FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;