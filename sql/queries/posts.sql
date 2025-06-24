-- name: AddPost :one
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    feed_id,
    title,
    url,
    description,
    published_at
) VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetPostsforUser :many
WITH follow AS (
    SELECT feed_id, user_id FROM feed_follows
    WHERE user_id = $1
)
SELECT posts.*
FROM follow INNER JOIN posts ON follow.feed_id = posts.feed_id
ORDER BY posts.published_at DESC
LIMIT $2;
