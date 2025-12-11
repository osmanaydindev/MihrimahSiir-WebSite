#!/bin/bash

# VPS Docker Deployment Script
# Bu script'i VPS'te çalıştıracaksınız

set -e

echo "🚀 VPS'e Docker deployment başlıyor..."

# 1. Veritabanını yedekle
echo "🗄️ Mevcut veritabanını yedekliyorum..."
./scripts/backup-database.sh

# 2. Uploads volume kontrol et ve oluştur
echo "📦 Uploads volume kontrol ediliyor..."
if ! docker volume inspect uploads_data >/dev/null 2>&1; then
    echo "📁 Yeni uploads volume oluşturuluyor..."
    docker volume create uploads_data

    # Eğer host'ta uploads klasörü varsa, volume'a kopyala
    if [ -d "./backend/uploads" ]; then
        echo "📋 Mevcut uploads dosyaları volume'a kopyalanıyor..."
        docker run --rm \
            -v $(pwd)/backend/uploads:/source:ro \
            -v uploads_data:/dest \
            alpine sh -c "cp -r /source/* /dest/ 2>/dev/null || echo 'Klasör boş veya kopyalanamadı'"
    fi
fi

# 3. Eski container'ları durdur
echo "🛑 Eski container'ları durduruyorum..."
docker-compose -f docker-compose.vps.yml down || true

# 3. Yeni image'ları build et
echo "🔨 Yeni image'ları build ediyorum..."
docker-compose -f docker-compose.vps.yml build --no-cache

# 4. Container'ları başlat
echo "▶️ Container'ları başlatıyorum..."
docker-compose -f docker-compose.vps.yml up -d

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
    echo "⚠️ Backend health check başarısız (normal olabilir)"
fi

# Frontend health check
if curl -f http://localhost:80 > /dev/null 2>&1; then
    echo "✅ Frontend sağlıklı"
else
    echo "❌ Frontend health check başarısız"
fi

echo "🎉 VPS deployment tamamlandı!"
echo "📊 Container durumu:"
docker-compose -f docker-compose.vps.yml ps

echo "🌐 Uygulama şu adreste çalışıyor:"
echo "   Frontend: http://your-domain.com"
echo "   Backend: http://your-domain.com/api"
