-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES(
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)

SELECT a.*, feeds.name AS feed_name, users.name AS user_name
FROM inserted_feed_follow a 
INNER JOIN users on users.id = a.user_id
INNER JOIN feeds on feeds.id = a.feed_id;