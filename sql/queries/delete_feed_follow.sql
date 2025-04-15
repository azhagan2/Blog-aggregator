-- name: Delete_Feed_Follow :exec

DELETE FROM feed_follows a
WHERE a.user_id = $1 and a.feed_id = $2;