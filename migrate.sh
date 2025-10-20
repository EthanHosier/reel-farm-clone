#!/bin/bash

# Migration runner for Reel Farm
# Usage: ./migrate.sh [up|down] [migration_number]

set -e

MIGRATIONS_DIR="migrations"
DATABASE_URL="${DATABASE_URL:-${DB_URL:-postgresql://localhost:5432/reel_farm}}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[MIGRATE]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Check if psql is available
if ! command -v psql &> /dev/null; then
    error "psql is required but not installed"
fi

# Create migrations table if it doesn't exist
init_migrations() {
    log "Initializing migrations table..."
    psql "$DATABASE_URL" -c "
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(255) PRIMARY KEY,
            applied_at TIMESTAMPTZ DEFAULT NOW()
        );
    "
}

# Get latest migration version
get_latest_version() {
    psql "$DATABASE_URL" -t -c "SELECT COALESCE(MAX(version), '000') FROM schema_migrations;" | tr -d ' '
}

# Apply migration
apply_migration() {
    local version=$1
    local file=$(find "$MIGRATIONS_DIR" -name "${version}_*.sql" | head -1)
    
    if [ -z "$file" ] || [ ! -f "$file" ]; then
        error "Migration file not found for version: $version"
    fi
    
    log "Applying migration $version..."
    psql "$DATABASE_URL" -f "$file"
    
    # Record migration
    psql "$DATABASE_URL" -c "INSERT INTO schema_migrations (version) VALUES ('$version') ON CONFLICT DO NOTHING;"
    
    log "Migration $version applied successfully"
}

# Run all pending migrations
migrate_up() {
    init_migrations
    
    local latest=$(get_latest_version)
    log "Latest applied migration: $latest"
    
    for file in $MIGRATIONS_DIR/*.sql; do
        if [ -f "$file" ]; then
            version=$(basename "$file" | cut -d'_' -f1)
            if [[ "$version" > "$latest" ]]; then
                apply_migration "$version"
            fi
        fi
    done
    
    log "All migrations applied"
}

# Show migration status
status() {
    init_migrations
    
    echo -e "${YELLOW}Migration Status:${NC}"
    echo "Applied migrations:"
    psql "$DATABASE_URL" -c "SELECT version, applied_at FROM schema_migrations ORDER BY version;"
    
    echo -e "\nAvailable migrations:"
    for file in $MIGRATIONS_DIR/*.sql; do
        if [ -f "$file" ]; then
            version=$(basename "$file" | cut -d'_' -f1)
            name=$(basename "$file" .sql | cut -d'_' -f2-)
            echo "  $version: $name"
        fi
    done
}

# Main
case "${1:-up}" in
    "up")
        migrate_up
        ;;
    "status")
        status
        ;;
    *)
        echo "Usage: $0 [up|status]"
        echo "  up     - Apply all pending migrations"
        echo "  status - Show migration status"
        exit 1
        ;;
esac
