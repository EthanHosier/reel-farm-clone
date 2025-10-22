-- name: GetUserAccount :one
SELECT * FROM public.user_accounts
WHERE id = $1;

-- name: GetUserByBillingCustomerID :one
SELECT * FROM public.user_accounts
WHERE billing_customer_id = $1;

-- name: UpdateUserPlan :exec
UPDATE public.user_accounts 
SET plan = $2, plan_started_at = $3, plan_ends_at = $4, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserBillingCustomerID :exec
UPDATE public.user_accounts 
SET billing_customer_id = $2, updated_at = NOW()
WHERE id = $1;

-- name: AddCreditsToUser :exec
UPDATE public.user_accounts 
SET credits = credits + $2, updated_at = NOW()
WHERE id = $1;

-- name: RemoveCreditsFromUser :exec
UPDATE public.user_accounts
SET credits = credits - $2, updated_at = NOW()
WHERE id = $1;
