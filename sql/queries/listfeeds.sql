-- name: GetFeeds :many
SELECT f.*, u.name AS user_name
FROM feeds AS f
JOIN users AS u ON f.user_id = u.id;
