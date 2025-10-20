-- name: GetUserAccount :one
SELECT * FROM public.user_accounts
WHERE id = $1;
