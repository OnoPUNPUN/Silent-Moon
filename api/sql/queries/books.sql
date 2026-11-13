-- name: CreateBook :one
INSERT INTO books (
    user_id,
    title,
    author,
    description,
    format,
    storage_key,
    cover_key,
    file_size,
    checksum,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetBookByID :one
SELECT * FROM books
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListBooksByUserID :many
SELECT * FROM books
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateBook :one
UPDATE books
SET title = $3,
    author = $4,
    description = $5,
    cover_key = $6,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateBookStatus :one
UPDATE books
SET status = $3,
    file_size = $4,
    checksum = $5,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteBook :exec
UPDATE books
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;
