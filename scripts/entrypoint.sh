#!/bin/sh
set -e

HOST="db"
USER="openconnect"
DB="openconnect"
export PGPASSWORD=1234

echo "Waiting for PostgreSQL..."
until pg_isready -h "$HOST" -U "$USER" -d "$DB"; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

echo "Running migrations..."
/usr/local/bin/migrate -path /migrations -database "postgres://openconnect:1234@db:5432/openconnect?sslmode=disable" up

echo "Starting application..."
exec ./main
