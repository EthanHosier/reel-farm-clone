-- Migration: Create credit transactions table
-- Description: Creates a table to track credit transactions for idempotency and audit purposes

-- Create the credit transactions table
CREATE TABLE public.credit_txns (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES public.user_accounts(id) ON DELETE CASCADE,
  request_id TEXT NOT NULL,
  amount INTEGER NOT NULL CHECK (amount > 0),
  status TEXT NOT NULL CHECK (status IN ('reserved','captured','refunded')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (request_id)
);

-- Add indexes for performance
CREATE INDEX idx_credit_txns_user_id ON public.credit_txns(user_id);
CREATE INDEX idx_credit_txns_status ON public.credit_txns(status);
CREATE INDEX idx_credit_txns_created_at ON public.credit_txns(created_at);

-- Add comments for documentation
COMMENT ON TABLE public.credit_txns IS 'Tracks credit transactions for idempotency and audit purposes';
COMMENT ON COLUMN public.credit_txns.id IS 'Unique transaction identifier';
COMMENT ON COLUMN public.credit_txns.user_id IS 'User who owns this transaction';
COMMENT ON COLUMN public.credit_txns.request_id IS 'Client-supplied idempotency key to prevent double-charging';
COMMENT ON COLUMN public.credit_txns.amount IS 'Number of credits involved in this transaction (must be > 0)';
COMMENT ON COLUMN public.credit_txns.status IS 'Transaction status: reserved (pending), captured (completed), refunded (failed)';
COMMENT ON COLUMN public.credit_txns.created_at IS 'When the transaction was first created';
COMMENT ON COLUMN public.credit_txns.updated_at IS 'When the transaction was last updated';
