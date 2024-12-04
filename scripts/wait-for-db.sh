#!/bin/sh
set -e

HOST="db"
USER="openconnect"
DB="openconnect"
export PGPASSWORD=1234

until pg_isready -h "$HOST" -U "$USER" -d "$DB"; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

echo "PostgreSQL is ready!"
exec "$@"