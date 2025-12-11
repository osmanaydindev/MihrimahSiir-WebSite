#!/bin/bash

# Veritabanı Restore Script'i
# Bu script Docker container başlatıldığında çalışır

set -e

echo "🔄 Veritabanı restore işlemi başlıyor..."

# Eğer backup dosyası varsa restore et
if [ -f "/backups/latest_backup.sql" ]; then
    echo "📥 Backup dosyası bulundu, restore ediliyor..."
    psql -U osmanaydin -d poem_blog_locale -f /backups/latest_backup.sql
    echo "✅ Veritabanı başarıyla restore edildi"
else
    echo "ℹ️ Backup dosyası bulunamadı, yeni veritabanı oluşturuluyor"
fi
