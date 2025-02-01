-- +goose Up

CREATE TABLE IF NOT EXISTS feed_follows(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users(id),
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id),
    UNIQUE(user_id,feed_id)
);

-- +goose Down
DROP TABLE IF EXISTS feed_follows;