#!/bin/bash

# Скрипт для автоматического применения миграций PostgreSQL
# Использует migrate CLI tool для управления миграциями

set -e

# Настройки подключения к базе данных
DB_HOST="${POSTGRES_HOST:-localhost}"
DB_PORT="${POSTGRES_PORT:-5432}"
DB_NAME="${POSTGRES_DB:-warmhouse}"
DB_USER="${POSTGRES_USER:-postgres}"
DB_PASSWORD="${POSTGRES_PASSWORD}"

# URL подключения к базе данных postgres для создания базы warmhouse
POSTGRES_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/postgres?sslmode=disable"
# URL подключения к базе данных warmhouse для миграций
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Директория с миграциями
MIGRATIONS_DIR="/migrations"

echo "Starting database migrations..."
echo "Database: ${DB_HOST}:${DB_PORT}/${DB_NAME}"
echo "User: ${DB_USER}"

# Ждем, пока база данных станет доступной
echo "Waiting for database to be ready..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is ready!"

# Создаем базу данных warmhouse, если она не существует
echo "Creating database warmhouse if it doesn't exist..."
psql "$POSTGRES_URL" -c "CREATE DATABASE warmhouse;" || echo "Database warmhouse already exists or creation failed"

# Проверяем, установлен ли migrate
if ! command -v migrate &> /dev/null; then
    echo "Installing migrate CLI tool..."
    # Устанавливаем migrate для Linux
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
    chmod +x migrate
    mv migrate /usr/local/bin/
fi

# Создаем таблицу для отслеживания миграций, если её нет
echo "Creating migrations table..."
migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version || true

# Применяем миграции
echo "Applying migrations from ${MIGRATIONS_DIR}..."
migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up

echo "Migrations completed successfully!"
