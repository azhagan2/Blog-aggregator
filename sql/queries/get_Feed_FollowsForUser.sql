-- name: GetFeedFollowsForUser :many
SELECT feeds.name
FROM feed_follows a 
INNER JOIN users on users.id = a.user_id
INNER JOIN feeds on feeds.id = a.feed_id
WHERE users.name = $1;