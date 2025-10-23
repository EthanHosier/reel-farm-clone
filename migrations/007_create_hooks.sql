-- Migration: Create hooks table
-- Description: Creates a table to store individual generated hooks for users

-- Create the hooks table
CREATE TABLE public.hooks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES public.user_accounts(id) ON DELETE CASCADE,
  generation_id UUID NOT NULL, -- Groups hooks from the same generation
  prompt TEXT NOT NULL,
  hook_text TEXT NOT NULL,
  hook_index INTEGER NOT NULL CHECK (hook_index >= 0), -- Order within the generation
  credits_used INTEGER NOT NULL CHECK (credits_used > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Add indexes for performance
CREATE INDEX idx_hooks_user_id ON public.hooks(user_id);
CREATE INDEX idx_hooks_generation_id ON public.hooks(generation_id);
CREATE INDEX idx_hooks_created_at ON public.hooks(created_at);

-- Add comments for documentation
COMMENT ON TABLE public.hooks IS 'Stores individual generated hooks for users';
COMMENT ON COLUMN public.hooks.id IS 'Unique hook identifier';
COMMENT ON COLUMN public.hooks.user_id IS 'User who generated this hook';
COMMENT ON COLUMN public.hooks.generation_id IS 'Groups hooks from the same generation request';
COMMENT ON COLUMN public.hooks.prompt IS 'The prompt used to generate this hook';
COMMENT ON COLUMN public.hooks.hook_text IS 'The actual hook text';
COMMENT ON COLUMN public.hooks.hook_index IS 'Order of this hook within the generation (0-based)';
COMMENT ON COLUMN public.hooks.credits_used IS 'Number of credits consumed for this generation';
COMMENT ON COLUMN public.hooks.created_at IS 'When the hook was generated';
COMMENT ON COLUMN public.hooks.updated_at IS 'When the record was last updated';
