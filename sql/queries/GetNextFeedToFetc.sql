-- name: Get_Next_Feed_to_fetch :one

SELECT *
FROM feeds 
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;