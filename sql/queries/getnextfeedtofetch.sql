-- name: GetNextFeedToFetch :one

SELECT * FROM feeds AS f
ORDER BY f.last_fetched_at ASC NULLS FIRST
LIMIT 1;
