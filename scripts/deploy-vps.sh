#!/bin/bash

# VPS Docker Deployment Script
# Bu script'i VPS'te çalıştıracaksınız

set -e

COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.vps.yml}"
if docker compose version >/dev/null 2>&1; then
    COMPOSE=(docker compose -f "$COMPOSE_FILE")
elif command -v docker-compose >/dev/null 2>&1; then
    COMPOSE=(docker-compose -f "$COMPOSE_FILE")
else
    echo "❌ docker compose bulunamadı"
    exit 1
fi

echo "🚀 VPS'e Docker deployment başlıyor..."

# 1. Veritabanını yedekle
echo "🗄️ Mevcut veritabanını yedekliyorum..."
./scripts/backup-database.sh

# 2. Uploads volume kontrol et ve oluştur
echo "📦 Uploads volume kontrol ediliyor..."
UPLOADS_VOLUME="${UPLOADS_VOLUME:-poem_uploads_data}"
if ! docker volume inspect "$UPLOADS_VOLUME" >/dev/null 2>&1; then
    echo "📁 Yeni uploads volume oluşturuluyor..."
    docker volume create "$UPLOADS_VOLUME"

    # Eğer host'ta uploads klasörü varsa, volume'a kopyala
    if [ -d "./backend/uploads" ]; then
        echo "📋 Mevcut uploads dosyaları volume'a kopyalanıyor..."
        docker run --rm \
            -v $(pwd)/backend/uploads:/source:ro \
            -v "$UPLOADS_VOLUME":/dest \
            alpine sh -c "cp -r /source/* /dest/ 2>/dev/null || echo 'Klasör boş veya kopyalanamadı'"
    fi
fi

# 3. Yeni image'ları build et
echo "🔨 Yeni image'ları build ediyorum..."
"${COMPOSE[@]}" build

# 4. Container'ları başlat
echo "▶️ Container'ları başlatıyorum..."
"${COMPOSE[@]}" up -d

# 5. Eski image'ları temizle (hafıza optimizasyonu)
echo "🧹 Eski image'lar temizleniyor..."
docker image prune -f

# 5. Health check
echo "🏥 Health check yapıyorum..."
sleep 15

# Backend health check
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Backend sağlıklı"
else
    echo "❌ Backend health check başarısız"
fi

# Frontend health check
if curl -f http://localhost:80 > /dev/null 2>&1; then
    echo "✅ Frontend sağlıklı"
else
    echo "❌ Frontend health check başarısız"
fi

echo "🎉 VPS deployment tamamlandı!"
echo "📊 Container durumu:"
"${COMPOSE[@]}" ps

echo "🌐 Uygulama şu adreste çalışıyor:"
echo "   Frontend: http://your-domain.com"
echo "   Backend: http://your-domain.com/api"
