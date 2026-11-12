-- name: CreateTag :one
INSERT INTO tags (
    user_id,
    name
) VALUES (
    $1, $2
)
ON CONFLICT (user_id, name) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: ListTagsByUserID :many
SELECT * FROM tags
WHERE user_id = $1
ORDER BY name ASC;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE id = $1 AND user_id = $2;

-- name: AddTagToBook :exec
INSERT INTO book_tags (
    book_id,
    tag_id
) VALUES (
    $1, $2
)
ON CONFLICT (book_id, tag_id) DO NOTHING;

-- name: RemoveTagFromBook :exec
DELETE FROM book_tags
WHERE book_id = $1 AND tag_id = $2;

-- name: ListTagsForBook :many
SELECT t.* FROM tags t
JOIN book_tags bt ON t.id = bt.tag_id
WHERE bt.book_id = $1 AND t.user_id = $2
ORDER BY t.name ASC;
