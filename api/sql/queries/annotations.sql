-- name: CreateAnnotation :one
INSERT INTO annotations (
    id,
    user_id,
    book_id,
    device_id,
    type,
    selected_text,
    color,
    note,
    location,
    version
) VALUES (
    COALESCE(sqlc.narg('id'), gen_random_uuid()),
    $1, $2, $3, $4, $5, $6, $7, $8, 1
)
RETURNING *;

-- name: GetAnnotationByID :one
SELECT * FROM annotations
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListAnnotationsByBookID :many
SELECT * FROM annotations
WHERE book_id = $1 AND user_id = $2 AND deleted_at IS NULL
ORDER BY created_at ASC;

-- name: UpdateAnnotation :one
UPDATE annotations
SET type = COALESCE(sqlc.narg('type'), type),
    selected_text = COALESCE(sqlc.narg('selected_text'), selected_text),
    color = COALESCE(sqlc.narg('color'), color),
    note = COALESCE(sqlc.narg('note'), note),
    location = COALESCE(sqlc.narg('location'), location),
    version = version + 1,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteAnnotation :exec
UPDATE annotations
SET deleted_at = NOW(),
    version = version + 1,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;
