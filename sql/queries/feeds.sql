-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name as username
FROM feeds 
INNER JOIN users
ON feeds.user_id = users.id;

-- name: GetFeedIDByURL :one
SELECT id FROM feeds WHERE url=$1;

-- name: MarkFeedFetched :exec
UPDATE feeds SET
last_fetched_at = $2,
updated_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM FEEDS ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;