## Supabase Configuration

1. **Enable Magic Link**: Magic Link is enabled by default in Supabase
2. **Set Site URL**: In your Supabase dashboard, go to Authentication > URL Configuration and set:
   - Site URL: `http://localhost:5173` (for development)
   - Redirect URLs: `http://localhost:5173/dashboard`

## How It Works

1. User enters email on home page (`/`)
2. Supabase sends magic link to their email
3. User clicks the link, which redirects directly to `/dashboard`
4. Supabase handles authentication automatically

## Routes

- `/` - Magic link login form (existing auth page)
- `/dashboard` - Protected dashboard (placeholder)

## Usage

1. Start your dev server: `npm run dev`
2. Navigate to `/` (home page)
3. Enter your email and click "Send Magic Link"
4. Check your email and click the magic link
5. You'll be redirected directly to the dashboard
