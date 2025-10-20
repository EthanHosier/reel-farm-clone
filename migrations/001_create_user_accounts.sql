-- Migration: 001_create_user_accounts.sql
-- Description: Create user_accounts table with plan management
-- Created: 2025-01-20

-- App-level table keyed to auth.users
create table public.user_accounts (
  id uuid primary key references auth.users(id) on delete cascade,
  plan text not null default 'free' check (plan in ('free','pro','enterprise')),
  plan_started_at timestamptz not null default now(),
  plan_ends_at timestamptz,
  billing_customer_id text,          -- e.g., Stripe customer id
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

-- updated_at maintenance
create or replace function public.tg_set_updated_at()
returns trigger language plpgsql as $$
begin
  new.updated_at := now();
  return new;
end;
$$;

create trigger set_updated_at_user_accounts
before update on public.user_accounts
for each row execute function public.tg_set_updated_at();

-- Auto-insert "free" row on signup
create or replace function public.handle_new_auth_user()
returns trigger language plpgsql security definer as $$
begin
  insert into public.user_accounts (id) values (new.id);
  return new;
end;
$$;

create trigger on_auth_user_created
after insert on auth.users
for each row execute function public.handle_new_auth_user();
