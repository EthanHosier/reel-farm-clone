CREATE TABLE user_generated_videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES user_accounts(id),
    ai_avatar_video_id UUID NOT NULL REFERENCES ai_avatar_videos(id),
    overlay_text TEXT NOT NULL,
    generated_video_filename VARCHAR(255) NOT NULL, -- S3 filename in user-generated-videos/ folder
    thumbnail_filename VARCHAR(255) NOT NULL, -- S3 filename in user-generated-videos/ folder
    status VARCHAR(20) DEFAULT 'processing', -- processing, completed, failed
    error_message TEXT, -- Store error details if processing fails
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add updated_at trigger
CREATE TRIGGER set_updated_at_user_generated_videos
BEFORE UPDATE ON public.user_generated_videos
FOR EACH ROW EXECUTE FUNCTION public.tg_set_updated_at();

-- Add indexes for performance
CREATE INDEX idx_user_generated_videos_user_id ON public.user_generated_videos(user_id);
CREATE INDEX idx_user_generated_videos_ai_avatar_video_id ON public.user_generated_videos(ai_avatar_video_id);
CREATE INDEX idx_user_generated_videos_status ON public.user_generated_videos(status);
CREATE INDEX idx_user_generated_videos_created_at ON public.user_generated_videos(created_at);
