-- name: GetUser :one
SELECT * FROM users
WHERE public_id = ? LIMIT 1;