-- name: GetFeedFollowsForUser :many
SELECT COUNT(*) AS follows, users.name AS UserName, feeds.name AS FeedName FROM feed_follows AS ff 
JOIN users ON users.id = ff.user_id
JOIN feeds ON feeds.id = ff.feed_id
WHERE users.name = $1
GROUP BY UserName, FeedName;

