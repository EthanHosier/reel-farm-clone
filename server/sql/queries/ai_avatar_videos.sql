-- name: GetAllVideos :many
SELECT * FROM ai_avatar_videos
ORDER BY created_at DESC;

-- name: GetVideoByID :one
SELECT * FROM ai_avatar_videos
WHERE id = $1;

-- name: CreateVideo :one
INSERT INTO ai_avatar_videos (
    id,
    title,
    description,
    filename,
    thumbnail_filename,
    duration,
    file_size
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdateVideo :one
UPDATE ai_avatar_videos
SET 
    title = $2,
    description = $3,
    duration = $4,
    file_size = $5
WHERE id = $1
RETURNING *;

-- name: DeleteVideo :exec
DELETE FROM ai_avatar_videos
WHERE id = $1;
