-- name: ReserveCredits :one
INSERT INTO public.credit_txns (user_id, request_id, amount, status)
VALUES ($1, $2, $3, 'reserved')
ON CONFLICT (request_id) DO UPDATE
  SET updated_at = NOW()
RETURNING id, status;

-- name: AtomicDebitCredits :one
UPDATE public.user_accounts
SET credits = credits - $2, updated_at = NOW()
WHERE id = $1 AND credits >= $2
RETURNING credits;

-- name: CaptureCredits :exec
UPDATE public.credit_txns
SET status = 'captured', updated_at = NOW()
WHERE id = $1 AND status = 'reserved';

-- name: RefundCredits :exec
UPDATE public.user_accounts
SET credits = credits + $2, updated_at = NOW()
WHERE id = $1;

-- name: MarkTxnRefunded :exec
UPDATE public.credit_txns
SET status = 'refunded', updated_at = NOW()
WHERE id = $1;

-- name: GetTxnStatus :one
SELECT status FROM public.credit_txns WHERE id = $1;

-- name: GetTxnByRequestID :one
SELECT id, user_id, request_id, amount, status, created_at, updated_at 
FROM public.credit_txns 
WHERE request_id = $1;

-- name: GetStaleReservedTxns :many
SELECT id, user_id, amount 
FROM public.credit_txns 
WHERE status = 'reserved' 
  AND created_at < NOW() - INTERVAL '10 minutes'
ORDER BY created_at ASC;
