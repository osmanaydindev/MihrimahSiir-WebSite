#!/bin/bash

# VPS Veritabanı Yedekleme Script'i
# Bu script'i VPS'te çalıştıracaksınız
# Environment variables: DB_USER, DB_NAME, DB_HOST (optional)

set -e

# Default values
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-your_database_name}"
DB_HOST="${DB_HOST:-localhost}"
BACKUP_DIR="${BACKUP_DIR:-/backups}"

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="poem_backup_${TIMESTAMP}.sql"

echo "🗄️ VPS veritabanını yedekliyorum..."
echo "Database: $DB_NAME, User: $DB_USER, Host: $DB_HOST"

# Backup dizinini oluştur
mkdir -p $BACKUP_DIR

# PostgreSQL veritabanını yedekle
pg_dump -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" > "${BACKUP_DIR}/${BACKUP_FILE}"

echo "✅ Yedekleme tamamlandı: ${BACKUP_DIR}/${BACKUP_FILE}"

# Son 5 yedeği tut, eskilerini sil
cd $BACKUP_DIR
ls -t poem_backup_*.sql | tail -n +6 | xargs -r rm

echo "📊 Mevcut yedekler:"
ls -la poem_backup_*.sql
