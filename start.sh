#!/bin/sh


set -e

echo "run db migratioon"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up
exec "$@"