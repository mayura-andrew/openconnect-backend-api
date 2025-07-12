#!/bin/sh
set -e

HOST="db"
USER="openconnect"
DB="openconnect"
export PGPASSWORD=1234

echo "Waiting for PostgreSQL..."
until pg_isready -h "$HOST" -U "$USER" -d "$DB"; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done

echo "PostgreSQL is up - executing command"

# Wait a bit more to ensure database is fully ready
sleep 5

echo "Running migrations..."
/usr/local/bin/migrate -path /migrations -database "postgres://openconnect:1234@db:5432/openconnect?sslmode=disable" up

if [ $? -eq 0 ]; then
    echo "Migrations completed successfully"
else
    echo "Migration failed with exit code $?"
    exit 1
fi

echo "Starting application..."
exec ./main