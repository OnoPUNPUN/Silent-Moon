-- name: CreateDevice :one
INSERT INTO devices (
    user_id,
    name,
    platform
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetDeviceByID :one
SELECT * FROM devices
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: ListDevicesByUserID :many
SELECT * FROM devices
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateDevice :one
UPDATE devices
SET name = $3
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: UpdateDeviceSyncTime :exec
UPDATE devices
SET last_synced_at = $3
WHERE id = $1 AND user_id = $2;

-- name: RevokeDevice :exec
UPDATE devices
SET revoked_at = NOW()
WHERE id = $1 AND user_id = $2;
