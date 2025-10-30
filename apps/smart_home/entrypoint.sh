#!/bin/sh
set -e

echo "Starting Smart Home application..."

# Apply database migrations if needed
if [ -n "$DATABASE_URL" ]; then
    echo "Applying database schema updates..."
    
    # Add device_id column if it doesn't exist
    psql "$DATABASE_URL" <<SQL || true
        ALTER TABLE sensors ADD COLUMN IF NOT EXISTS device_id TEXT NOT NULL DEFAULT '';
        CREATE INDEX IF NOT EXISTS idx_sensors_device_id ON sensors(device_id);
SQL
    
    echo "Database schema updated successfully!"
fi

# Start the application
exec "$@"

