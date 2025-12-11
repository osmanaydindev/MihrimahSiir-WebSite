#!/bin/bash

# Docker Deployment Script
# Bu script'i VPS'te çalıştıracaksınız

set -e

echo "🚀 Starting Docker deployment..."

# Environment kontrolü
if [ ! -f "env.production" ]; then
    echo "❌ env.production file not found!"
    exit 1
fi

# Eski container'ları durdur ve sil
echo "🛑 Stopping existing containers..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml down

# Image'ları yeniden build et
echo "🔨 Building new images..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build --no-cache

# Container'ları başlat
echo "▶️ Starting containers..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Health check
echo "🏥 Checking health..."
sleep 10

# Backend health check
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Backend is healthy"
else
    echo "❌ Backend health check failed"
fi

# Frontend health check
if curl -f http://localhost:3000 > /dev/null 2>&1; then
    echo "✅ Frontend is healthy"
else
    echo "❌ Frontend health check failed"
fi

echo "🎉 Deployment completed!"
echo "📊 Container status:"
docker-compose -f docker-compose.yml -f docker-compose.prod.yml ps
