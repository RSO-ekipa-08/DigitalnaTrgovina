-- name: GetApplication :one
SELECT 
    id, name, description, developer_id, category, price, size,
    min_android_version, current_version, tags, screenshots, storage_url,
    rating, downloads, created_at, updated_at
FROM applications
WHERE id = $1 LIMIT 1;

-- name: ListApplications :many
SELECT 
    id, name, description, developer_id, category, price, size,
    min_android_version, current_version, tags, screenshots, storage_url,
    rating, downloads, created_at, updated_at
FROM applications
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchApplications :many
SELECT 
    id, name, description, developer_id, category, price, size,
    min_android_version, current_version, tags, screenshots, storage_url,
    rating, downloads, created_at, updated_at
FROM applications
WHERE
    CASE
        WHEN @query::text != '' THEN
            search_vector @@ plainto_tsquery('english', @query::text)
        ELSE true
    END
    AND CASE
        WHEN @category::text != '' THEN
            category = @category::text
        ELSE true
    END
    AND CASE
        WHEN @min_price::decimal >= 0 THEN
            price >= @min_price::decimal
        ELSE true
    END
    AND CASE
        WHEN @max_price::decimal >= 0 THEN
            price <= @max_price::decimal
        ELSE true
    END
    AND CASE
        WHEN @min_android_version::text != '' THEN
            min_android_version >= @min_android_version::text
        ELSE true
    END
    AND CASE
        WHEN array_length(@tags::text[], 1) > 0 THEN
            tags && @tags::text[]
        ELSE true
    END
ORDER BY
    CASE
        WHEN @sort_by_downloads::boolean = true THEN downloads
        ELSE NULL
    END DESC NULLS LAST,
    CASE
        WHEN @sort_by_rating::boolean = true THEN rating
        ELSE NULL
    END DESC NULLS LAST,
    created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateApplication :one
INSERT INTO applications (
    name,
    description,
    developer_id,
    category,
    price,
    size,
    min_android_version,
    current_version,
    tags,
    screenshots,
    storage_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING 
    id, name, description, developer_id, category, price, size,
    min_android_version, current_version, tags, screenshots, storage_url,
    rating, downloads, created_at, updated_at;

-- name: UpdateApplication :one
UPDATE applications
SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    category = COALESCE(sqlc.narg('category'), category),
    price = COALESCE(sqlc.narg('price'), price),
    min_android_version = COALESCE(sqlc.narg('min_android_version'), min_android_version),
    current_version = COALESCE(sqlc.narg('current_version'), current_version),
    tags = COALESCE(sqlc.narg('tags'), tags),
    screenshots = COALESCE(sqlc.narg('screenshots'), screenshots),
    storage_url = COALESCE(sqlc.narg('storage_url'), storage_url),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING 
    id, name, description, developer_id, category, price, size,
    min_android_version, current_version, tags, screenshots, storage_url,
    rating, downloads, created_at, updated_at;

-- name: DeleteApplication :exec
DELETE FROM applications
WHERE id = $1;

-- name: IncrementDownloads :one
UPDATE applications
SET downloads = downloads + 1
WHERE id = $1
RETURNING 
    id, name, description, developer_id, category, price, size,
    min_android_version, current_version, tags, screenshots, storage_url,
    rating, downloads, created_at, updated_at;

-- name: ListCategories :many
SELECT id, name, description, created_at, updated_at
FROM categories
ORDER BY name ASC
LIMIT $1 OFFSET $2; 