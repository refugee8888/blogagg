-- name: GetPostsForUser :many
SELECT p.*, f.name AS feed_name
FROM posts AS p
JOIN feed_follows AS ff ON ff.feed_id = p.feed_id
JOIN feeds AS f ON f.id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC NULLS LAST
LIMIT $2;
