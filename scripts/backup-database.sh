#!/bin/bash

# VPS Veritabanı Yedekleme Script'i
# Environment variables:
# DB_CONTAINER, DB_USER, DB_NAME, BACKUP_DIR

set -e

# Production defaults for this project. BACKUP_DIR intentionally points to the
# deploy user's home, not /backups, so the script does not require root.
DB_CONTAINER="${DB_CONTAINER:-poem_database_vps}"
DB_USER="${DB_USER:-osmanaydin}"
DB_NAME="${DB_NAME:-poem_blog_locale}"
BACKUP_DIR="${BACKUP_DIR:-$HOME/poem_backups}"

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="poem_backup_${TIMESTAMP}.sql"

echo "🗄️ VPS veritabanını yedekliyorum..."
echo "Database: $DB_NAME, User: $DB_USER, Container: $DB_CONTAINER"

# Backup dizinini oluştur
mkdir -p "$BACKUP_DIR"

# PostgreSQL veritabanını container içinden yedekle. Böylece host'ta pg_dump
# kurulumu veya database portunun dışarı açılması gerekmez.
docker exec "$DB_CONTAINER" pg_dump -U "$DB_USER" -d "$DB_NAME" > "${BACKUP_DIR}/${BACKUP_FILE}"

echo "✅ Yedekleme tamamlandı: ${BACKUP_DIR}/${BACKUP_FILE}"

# Son 5 yedeği tut, eskilerini sil
cd "$BACKUP_DIR"
ls -t poem_backup_*.sql | tail -n +6 | xargs -r rm

echo "📊 Mevcut yedekler:"
ls -la poem_backup_*.sql
