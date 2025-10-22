-- Migration: Update credits column default value
-- Description: Changes the default value of credits column from 0 to 100

-- Update the default value for the credits column
ALTER TABLE public.user_accounts 
ALTER COLUMN credits SET DEFAULT 100;

-- Update existing users to have 100 credits (if they currently have 0)
UPDATE public.user_accounts 
SET credits = 100 
WHERE credits = 0;

-- Add comment to document the column
COMMENT ON COLUMN public.user_accounts.credits IS 'Number of credits available to the user (must be >= 0, default 100)';
