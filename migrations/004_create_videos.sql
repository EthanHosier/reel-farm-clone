-- Migration: 004_create_videos.sql
-- Description: Create videos table for AI avatar video content
-- Created: 2025-01-20

-- Videos table for AI avatar content
create table public.videos (
  id uuid primary key default gen_random_uuid(),
  title text not null,
  description text,
  filename text not null,                    -- e.g., "video-1.mp4"
  thumbnail_filename text not null,          -- e.g., "video-1.jpg"
  duration integer,                          -- duration in seconds
  file_size bigint,                          -- file size in bytes
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

-- Add updated_at trigger
create trigger set_updated_at_videos
before update on public.videos
for each row execute function public.tg_set_updated_at();

-- Add indexes for common queries
create index idx_videos_created_at on public.videos(created_at desc);
create index idx_videos_title on public.videos(title);
