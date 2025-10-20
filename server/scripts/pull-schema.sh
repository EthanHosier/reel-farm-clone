#!/bin/bash

# Script to pull schema from Supabase and generate SQLC code

set -e

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL environment variable is required"
    echo "Example: DATABASE_URL=postgresql://postgres:password@db.project.supabase.co:5432/postgres"
    exit 1
fi

echo "🔄 Pulling schema from database..."

# Create schema directory if it doesn't exist
mkdir -p sql/schema

# Pull schema from database (excluding auth schema and system tables)
# Use PostgreSQL 17 version to match Supabase
/opt/homebrew/opt/postgresql@17/bin/pg_dump "$DATABASE_URL" \
  --schema-only \
  --no-owner \
  --no-privileges \
  --no-sync \
  --exclude-schema=auth \
  --exclude-schema=extensions \
  --exclude-schema=graphql \
  --exclude-schema=graphql_public \
  --exclude-schema=realtime \
  --exclude-schema=storage \
  --exclude-schema=supabase_functions \
  --exclude-schema=supabase_migrations \
  --exclude-schema=vault \
  --exclude-table-data=* \
  --file=sql/schema/schema.sql

echo "✅ Schema pulled successfully"

# Clean up schema file (remove pg_dump artifacts)
echo "🧹 Cleaning schema file..."
sed -i '' '/^\\restrict/d' sql/schema/schema.sql
sed -i '' '/^\\unrestrict/d' sql/schema/schema.sql
sed -i '' '/^\\connect/d' sql/schema/schema.sql
sed -i '' '/^SET /d' sql/schema/schema.sql
sed -i '' '/^SELECT pg_catalog/d' sql/schema/schema.sql
sed -i '' '/^-- PostgreSQL database dump/d' sql/schema/schema.sql

# Generate SQLC code
echo "🔄 Generating SQLC code..."
~/go/bin/sqlc generate

echo "✅ SQLC code generated successfully"
echo "📁 Generated files are in the 'db' directory"
