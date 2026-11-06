-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = Now(), updated_at = Now()
WHERE id = $1
RETURNING *;



