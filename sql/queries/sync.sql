-- name: CreateSyncChange :one
INSERT INTO sync_changes (
    user_id,
    device_id,
    entity_type,
    entity_id,
    operation
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListSyncChangesAfterCursor :many
SELECT * FROM sync_changes
WHERE user_id = $1 AND id > $2
ORDER BY id ASC
LIMIT $3;

-- name: GetLatestSyncSequenceForUser :one
SELECT COALESCE(MAX(id), 0)::BIGINT AS latest_seq
FROM sync_changes
WHERE user_id = $1;

-- name: RecordSyncOperation :one
INSERT INTO sync_operations (
    operation_id,
    user_id,
    device_id,
    entity_type,
    entity_id,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetSyncOperationByID :one
SELECT * FROM sync_operations
WHERE operation_id = $1 AND user_id = $2
LIMIT 1;
