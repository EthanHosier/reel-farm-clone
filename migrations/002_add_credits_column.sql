-- Migration: Add credits column to user_accounts table
-- Description: Adds a credits column to track user credits (integer >= 0)

-- Add credits column with default value of 100 and check constraint
ALTER TABLE public.user_accounts 
ADD COLUMN credits INTEGER NOT NULL DEFAULT 0;

-- Add check constraint to ensure credits is >= 0
ALTER TABLE public.user_accounts 
ADD CONSTRAINT user_accounts_credits_check 
CHECK (credits >= 0);

-- Add comment to document the column
COMMENT ON COLUMN public.user_accounts.credits IS 'Number of credits available to the user (must be >= 0)';
