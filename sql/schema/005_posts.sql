-- +goose up

CREATE TABLE posts(
    Id INTEGER UNIQUE PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT UNIQUE NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP,
    feed_id INTEGER REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;