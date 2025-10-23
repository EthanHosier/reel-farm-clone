-- Migration: 005_rename_videos_to_ai_avatar_videos.sql
-- Description: Rename videos table to ai_avatar_videos
-- Created: 2025-01-20

-- Rename the videos table to ai_avatar_videos
ALTER TABLE public.videos RENAME TO ai_avatar_videos;

-- Rename the trigger
ALTER TRIGGER set_updated_at_videos ON public.ai_avatar_videos RENAME TO set_updated_at_ai_avatar_videos;

-- Rename the indexes
ALTER INDEX idx_videos_created_at RENAME TO idx_ai_avatar_videos_created_at;
ALTER INDEX idx_videos_title RENAME TO idx_ai_avatar_videos_title;
