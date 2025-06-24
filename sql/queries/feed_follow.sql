-- name: AddFeedFollow :one
WITH added_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2
    )
RETURNING *
)
SELECT added_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
    FROM added_feed_follow
INNER JOIN users ON added_feed_follow.user_id = users.id
INNER JOIN feeds ON added_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 and feed_id = $2;
