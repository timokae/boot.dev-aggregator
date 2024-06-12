-- name: CreateFeedFollow :one
INSERT INTO feed_follows (
  id,
  user_id,
  feed_id,
  created_at,
  updated_at
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollow :one
SELECT * FROM feed_follows WHERE id = $1 LIMIT 1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE id = $1;

-- name: GetFeedFollowsOfUser :many
SELECT * FROM feed_follows WHERE user_id = $1;