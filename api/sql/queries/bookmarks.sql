-- name: CreateBookmark :one
INSERT INTO bookmarks (
    id,
    user_id,
    book_id,
    device_id,
    title,
    location,
    version
) VALUES (
    COALESCE(sqlc.narg('id'), gen_random_uuid()),
    $1, $2, $3, $4, $5, 1
)
RETURNING *;

-- name: GetBookmarkByID :one
SELECT * FROM bookmarks
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListBookmarksByBookID :many
SELECT * FROM bookmarks
WHERE book_id = $1 AND user_id = $2 AND deleted_at IS NULL
ORDER BY created_at ASC;

-- name: UpdateBookmark :one
UPDATE bookmarks
SET title = COALESCE(sqlc.narg('title'), title),
    location = COALESCE(sqlc.narg('location'), location),
    version = version + 1,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteBookmark :exec
UPDATE bookmarks
SET deleted_at = NOW(),
    version = version + 1,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;
