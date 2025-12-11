#!/bin/bash

# Docker Hafıza Temizleme Script
# VPS'teki Docker hafıza sorunlarını çözer

set -e

echo "🧹 Docker hafıza temizliği başlıyor..."
echo ""

# Önce mevcut durumu göster
echo "📊 ŞU ANKİ DURUM:"
echo "=================="
docker system df
echo ""

# Kullanılmayan image'ları temizle
echo "🗑️  Kullanılmayan image'lar temizleniyor..."
docker image prune -f

echo ""
echo "🗑️  Durmuş container'lar temizleniyor..."
docker container prune -f

echo ""
echo "🗑️  Kullanılmayan network'ler temizleniyor..."
docker network prune -f

echo ""
echo "🗑️  Build cache temizleniyor..."
docker builder prune -f

echo ""
echo "✅ Temizlik tamamlandı!"
echo ""

# Sonraki durumu göster
echo "📊 YENİ DURUM:"
echo "=================="
docker system df
echo ""

echo "💡 İPUCU: Eğer hala çok fazla yer kullanılıyorsa:"
echo "   docker system prune -a -f  (Tüm kullanılmayan image'ları siler)"
echo ""
echo "⚠️  DİKKAT: Volume'ları silmek için asla --volumes kullanmayın!"
echo "   (Profil resimleri ve database kaybolur)"
