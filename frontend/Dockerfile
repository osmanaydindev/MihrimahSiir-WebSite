# Frontend Dockerfile
FROM node:18-alpine as builder

# Build argument for API URL
ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

WORKDIR /app

# Package dosyalarını kopyala
COPY package*.json ./

# Dependencies yükle (npm kullan)
RUN npm install

# Kaynak kodları kopyala
COPY . .

# Build et
RUN npm run build

# Production stage - Nginx ile serve et
FROM nginx:alpine

# Nginx config dosyasını kopyala
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Build edilmiş dosyaları kopyala
COPY --from=builder /app/dist /usr/share/nginx/html

# Port 3000'i expose et
EXPOSE 3000

# Nginx'i başlat
CMD ["nginx", "-g", "daemon off;"]
