-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
) 
RETURNING *;

-- name: GetPostsByUserID :many
SELECT * 
FROM posts 
INNER JOIN feed_follow ON posts.feed_id = feed_follow.feed_id
INNER JOIN users ON feed_follow.user_id = users.id
WHERE users.id = $1
ORDER BY posts.published_at DESC
LIMIT $2;