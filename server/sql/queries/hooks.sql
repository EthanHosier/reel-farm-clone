-- name: CreateHook :one
INSERT INTO public.hooks (user_id, generation_id, prompt, hook_text, hook_index, credits_used)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateHooksBatch :many
INSERT INTO public.hooks (user_id, generation_id, prompt, hook_text, hook_index, credits_used)
SELECT $1, $2, $3, unnest($4::text[]), unnest($5::int[]), $6
RETURNING *;

-- name: GetHooksByUser :many
SELECT * FROM public.hooks
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetHooksByGeneration :many
SELECT * FROM public.hooks
WHERE generation_id = $1
ORDER BY hook_index ASC;

-- name: GetHookByID :one
SELECT * FROM public.hooks
WHERE id = $1;

-- name: DeleteHook :exec
DELETE FROM public.hooks
WHERE id = $1 AND user_id = $2;

-- name: DeleteHooks :many
-- sqlc:arg hook_ids uuid[]
-- sqlc:arg user_id uuid
DELETE FROM public.hooks
WHERE id = ANY(@hook_ids::uuid[]) AND user_id = @user_id
RETURNING *;

-- name: GetUserHookCount :one
SELECT COUNT(*) FROM public.hooks
WHERE user_id = $1;
