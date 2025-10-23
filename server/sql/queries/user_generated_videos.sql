-- name: CreateUserGeneratedVideo :one
INSERT INTO user_generated_videos (
    id,
    user_id,
    ai_avatar_video_id,
    overlay_text,
    generated_video_filename,
    thumbnail_filename,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserGeneratedVideoByID :one
SELECT * FROM user_generated_videos WHERE id = $1;

-- name: GetUserGeneratedVideosByUserID :many
SELECT * FROM user_generated_videos WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateUserGeneratedVideoStatus :one
UPDATE user_generated_videos 
SET status = $2, error_message = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserGeneratedVideoFilenames :one
UPDATE user_generated_videos 
SET generated_video_filename = $2, thumbnail_filename = $3, status = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;
