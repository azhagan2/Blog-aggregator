-- +goose up
CREATE TABLE feed_follows(
    id INTEGER PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    feed_id INTEGER REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;