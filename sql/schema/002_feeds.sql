-- +goose up
CREATE TABLE feeds(
    id INTEGER PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;