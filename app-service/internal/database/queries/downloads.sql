-- name: CreateDownload :one
INSERT INTO downloads (
    user_id,
    application_id,
    ip_address,
    success,
    error_message
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetDownloadsByUser :many
SELECT * FROM downloads
WHERE user_id = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3;

-- name: GetDownloadsByApplication :many
SELECT * FROM downloads
WHERE application_id = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3;

-- name: GetDownloadStats :one
SELECT 
    COUNT(*) as total_downloads,
    COUNT(CASE WHEN success = true THEN 1 END) as successful_downloads,
    COUNT(CASE WHEN success = false THEN 1 END) as failed_downloads
FROM downloads
WHERE application_id = $1; 