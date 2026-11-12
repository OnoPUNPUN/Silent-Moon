-- name: CreateCollection :one
INSERT INTO collections (
    user_id,
    name,
    description
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetCollectionByID :one
SELECT * FROM collections
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListCollectionsByUserID :many
SELECT * FROM collections
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateCollection :one
UPDATE collections
SET name = $3,
    description = $4,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteCollection :exec
UPDATE collections
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: AddBookToCollection :exec
INSERT INTO book_collections (
    book_id,
    collection_id
) VALUES (
    $1, $2
)
ON CONFLICT (book_id, collection_id) DO NOTHING;

-- name: RemoveBookFromCollection :exec
DELETE FROM book_collections
WHERE book_id = $1 AND collection_id = $2;

-- name: ListBooksInCollection :many
SELECT b.* FROM books b
JOIN book_collections bc ON b.id = bc.book_id
WHERE bc.collection_id = $1
  AND b.user_id = $2
  AND b.deleted_at IS NULL
ORDER BY b.created_at DESC;
