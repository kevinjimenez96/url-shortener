-- name: GetMapping :one
SELECT * FROM mappings
WHERE path = $1;