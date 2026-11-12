-- name: UpsertReadingProgress :one
INSERT INTO reading_progress (
    user_id,
    book_id,
    device_id,
    location,
    progress_percentage,
    last_read_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW()
)
ON CONFLICT (book_id, device_id)
DO UPDATE SET
    location = EXCLUDED.location,
    progress_percentage = EXCLUDED.progress_percentage,
    last_read_at = EXCLUDED.last_read_at,
    updated_at = NOW()
RETURNING *;

-- name: GetLatestReadingProgressForBook :one
SELECT * FROM reading_progress
WHERE book_id = $1 AND user_id = $2
ORDER BY last_read_at DESC
LIMIT 1;

-- name: GetReadingProgressByDevice :one
SELECT * FROM reading_progress
WHERE book_id = $1 AND device_id = $2 AND user_id = $3
LIMIT 1;

-- name: ListReadingProgressForBook :many
SELECT * FROM reading_progress
WHERE book_id = $1 AND user_id = $2
ORDER BY last_read_at DESC;
