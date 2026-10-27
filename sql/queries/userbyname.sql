-- name: GetUser :one
SELECT * FROM users as u
WHERE u.name = $1;

