-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_id = (SELECT id FROM feeds WHERE url = $1)
AND feed_follows.user_id = $2;
