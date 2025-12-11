#!/bin/bash

# Sadece Backend Deployment Script
# Hafıza optimizasyonu ile backend'i günceller

set -e

echo "🚀 Backend deployment başlıyor..."

# 1. Git pull
echo "📥 Git pull yapılıyor..."
git pull origin main

# 2. Sadece backend'i yeniden build et ve başlat
echo "🔨 Backend build ediliyor..."
docker compose -f docker-compose.vps.yml up -d --build backend

# 3. Kullanılmayan image'ları temizle
echo "🧹 Eski image'lar temizleniyor..."
docker image prune -f

# 4. Durum kontrolü
echo "📊 Container durumları:"
docker compose -f docker-compose.vps.yml ps

# 5. Backend loglarını göster (son 20 satır)
echo "📜 Backend logları:"
docker compose -f docker-compose.vps.yml logs --tail=20 backend

echo "✅ Backend deployment tamamlandı!"
