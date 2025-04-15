-- name: GetPostsForUser :many
SELECT *
FROM posts 
JOIN feed_follows a ON posts.feed_id = a.feed_id  
JOIN feeds b ON a.feed_id = b.id
WHERE a.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;