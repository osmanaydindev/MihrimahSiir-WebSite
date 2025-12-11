# MihrimahSiir - Şiir ve Edebiyat Platformu

Modern, full-stack şiir ve edebiyat paylaşım platformu. Kullanıcılar şiirleri keşfedebilir, beğenebilir, yorum yapabilir ve arkadaşlarıyla paylaşabilir.

## Proje Özeti

MihrimahSiir, edebiyat severler için geliştirilmiş sosyal bir platformdur. Kullanıcılar şiirleri okuyabilir, favorilerine ekleyebilir, arkadaşlarıyla etkileşime geçebilir ve kendi profillerini yönetebilirler. Admin paneli ile içerik yönetimi kolayca yapılabilir.

## Teknoloji Yığını

### Backend
- **Go 1.23** - Performanslı backend geliştirme
- **Fiber v2** - Express benzeri, hızlı web framework
- **PostgreSQL** - Ana veritabanı
- **GORM** - ORM kütüphanesi
- **JWT** - Kullanıcı kimlik doğrulama
- **WebSocket** - Gerçek zamanlı iletişim
- **bcrypt** - Şifre hashleme

### Frontend
- **Vue 3** - Modern JavaScript framework
- **TypeScript** - Tip güvenliği
- **Vuetify 3** - Material Design UI kütüphanesi
- **Pinia** - State management
- **Vite** - Hızlı build aracı
- **Vue Router** - SPA routing
- **Axios** - HTTP client
- **Quill Editor** - Rich text editing
- **DOMPurify** - XSS koruması
- **VeeValidate** - Form validation

### DevOps & Araçlar
- **Docker & Docker Compose** - Konteynerizasyon
- **Nginx** - Reverse proxy ve static file serving
- **ESLint** - Code linting
- **Git** - Versiyon kontrolü

## Ana Özellikler

### Kullanıcı Özellikleri
- Kullanıcı kaydı ve girişi (JWT authentication)
- Profil yönetimi
- Şiirleri görüntüleme, arama ve filtreleme
- Şiirleri beğenme ve bookmark'lama
- Yorum yapma sistemi
- Arkadaşlık sistemi (istek gönder/kabul et/reddet)
- Okunan kitapları takip etme
- Hatırlatıcılar oluşturma
- En popüler ve en yeni şiirleri görüntüleme

### Admin Özellikleri
- Admin paneli
- Şiir CRUD işlemleri
- Kitap yönetimi
- Yazar yönetimi
- Kullanıcı yönetimi
- Anasayfa içerik yönetimi
- Mihrimah Card yönetimi
- Log görüntüleme ve yönetimi
- Yorum moderasyonu

### Güvenlik Özellikleri
- JWT tabanlı kimlik doğrulama
- Bcrypt ile şifre hashleme
- Rate limiting (global ve endpoint bazlı)
- Security headers (XSS, CSRF koruması)
- Input sanitization
- CORS konfigürasyonu
- File upload validation
- SQL injection koruması (GORM parametrized queries)

### Teknik Özellikler
- WebSocket ile gerçek zamanlı arkadaşlık bildirimleri
- Responsive tasarım
- SEO friendly routing
- Image upload ve optimize
- Lazy loading
- Pagination
- Advanced search
- Error handling ve logging
- Multi-stage Docker builds
- Health check endpoints

## Proje Yapısı

```
MihrimahSiir-WebSite/
├── backend/                    # Go Backend
│   ├── controllers/           # HTTP request handlers
│   ├── database/              # Database connection
│   ├── helpers/               # Helper functions (pagination, etc.)
│   ├── middlewares/           # Middleware (auth, rate limit, security)
│   ├── models/                # Database models
│   ├── routes/                # Route definitions
│   ├── security/              # Security utilities (sanitizer, validator)
│   ├── util/                  # Utilities (JWT)
│   ├── websocket/             # WebSocket hub and handlers
│   ├── main.go               # Application entry point
│   ├── go.mod                # Go dependencies
│   └── Dockerfile            # Backend Docker image
│
├── frontend/                   # Vue 3 Frontend
│   ├── src/
│   │   ├── components/       # Vue components
│   │   │   ├── common/      # Shared components
│   │   │   ├── frontend/    # User-facing components
│   │   │   ├── panel/       # Admin panel components
│   │   │   ├── profile/     # Profile components
│   │   │   └── friends/     # Friends system components
│   │   ├── composables/      # Vue 3 composables
│   │   ├── constants/        # Constants and configs
│   │   ├── layouts/          # Layout components
│   │   ├── pages/            # Page components
│   │   │   ├── Management/  # Admin pages
│   │   │   └── User/        # User pages
│   │   ├── plugins/          # Vue plugins (Vuetify, etc.)
│   │   ├── services/         # API services
│   │   │   └── api/         # API client and services
│   │   ├── stores/           # Pinia stores
│   │   ├── styles/           # Global styles
│   │   ├── types/            # TypeScript type definitions
│   │   ├── utils/            # Utility functions
│   │   ├── App.vue          # Root component
│   │   └── main.ts          # Application entry point
│   ├── public/               # Static assets
│   ├── nginx.conf           # Nginx configuration
│   ├── package.json         # Node dependencies
│   ├── tsconfig.json        # TypeScript configuration
│   ├── vite.config.ts       # Vite configuration
│   └── Dockerfile           # Frontend Docker image
│
├── scripts/                    # Deployment and utility scripts
│   ├── backup-database.sh
│   ├── cleanup-docker.sh
│   ├── deploy.sh
│   ├── deploy-backend-only.sh
│   ├── deploy-vps.sh
│   └── restore-database.sh
│
├── docker-compose.example.yml  # Docker Compose template
├── .env.docker.example        # Docker environment template
├── release.sh                 # Version release script
└── README.md                  # Bu dosya
```

## Kurulum

### Gereksinimler

- **Node.js** 18+ ve npm
- **Go** 1.23+
- **PostgreSQL** 15+
- **Docker & Docker Compose** (opsiyonel, önerilen)

### Yerel Geliştirme Ortamı

#### 1. Repository'yi Klonlayın

```bash
git clone <repository-url>
cd MihrimahSiir-WebSite
```

#### 2. Backend Kurulumu

```bash
cd backend

# Environment dosyasını oluşturun
cp .env.example .env

# .env dosyasını düzenleyin ve veritabanı bilgilerinizi girin
# DB_HOST=localhost
# DB_USER=your_user
# DB_PASSWORD=your_password
# DB_NAME=poem_blog
# DB_PORT=5432
# ADMIN_USERNAME=admin
# ADMIN_PASSWORD=your_admin_password

# Go modüllerini indirin
go mod download

# Uygulamayı çalıştırın
go run main.go
```

Backend varsayılan olarak `http://localhost:8080` adresinde çalışacaktır.

#### 3. Frontend Kurulumu

```bash
cd frontend

# Bağımlılıkları yükleyin
npm install

# Development server'ı başlatın
npm run dev
```

Frontend varsayılan olarak `http://localhost:3000` adresinde çalışacaktır.

### Docker ile Kurulum (Önerilen)

#### 1. Environment Dosyalarını Hazırlayın

```bash
# Docker Compose environment dosyasını oluşturun
cp .env.docker.example .env.docker

# Backend environment dosyasını oluşturun
cp backend/.env.example backend/.env
```

#### 2. Environment Dosyalarını Düzenleyin

`.env.docker` dosyası:
```env
DB_NAME=your_database_name
DB_USER=your_database_user
DB_PASSWORD=your_strong_password
```

`backend/.env` dosyası:
```env
DB_HOST=database
DB_USER=your_database_user
DB_PASSWORD=your_strong_password
DB_NAME=your_database_name
DB_PORT=5432
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_admin_password
APP_PORT=8080
CORS_ORIGINS=http://localhost:3000
```

#### 3. Docker Compose Dosyasını Oluşturun

```bash
cp docker-compose.example.yml docker-compose.yml
```

#### 4. Uygulamayı Başlatın

```bash
# Container'ları build edin ve başlatın
docker-compose up -d

# Logları görüntüleyin
docker-compose logs -f

# Container durumunu kontrol edin
docker-compose ps
```

Uygulama şu adreslerde erişilebilir olacaktır:
- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080`
- PostgreSQL: `localhost:5432`

#### 5. Container'ları Durdurma

```bash
# Container'ları durdur
docker-compose down

# Container'ları durdur ve volume'ları sil (VERİTABANI SİLİNİR!)
docker-compose down -v
```

## API Endpoints

### Kimlik Doğrulama
- `POST /register` - Yeni kullanıcı kaydı
- `POST /login` - Kullanıcı girişi
- `GET /auth-check` - Kimlik doğrulama kontrolü (protected)

### Şiirler
- `GET /poems` - Tüm şiirleri listele
- `GET /poem/:id` - Tek bir şiiri getir
- `POST /poem` - Yeni şiir oluştur (admin)
- `PUT /poem/:id` - Şiir güncelle (admin)
- `DELETE /poem/:id` - Şiir sil (admin)
- `GET /latest-poems` - En yeni şiirler
- `GET /popular-poems` - En popüler şiirler

### Kitaplar
- `GET /books` - Tüm kitapları listele
- `GET /book/:id` - Tek bir kitabı getir
- `POST /book` - Yeni kitap oluştur (admin)
- `PUT /book/:id` - Kitap güncelle (admin)
- `DELETE /book/:id` - Kitap sil (admin)

### Yazarlar
- `GET /authors` - Tüm yazarları listele
- `GET /author/:id` - Tek bir yazarı getir
- `POST /author` - Yeni yazar oluştur (admin)
- `PUT /author/:id` - Yazar güncelle (admin)
- `DELETE /author/:id` - Yazar sil (admin)

### Arkadaşlık Sistemi
- `POST /send-friend-request` - Arkadaşlık isteği gönder
- `GET /get-friend-requests` - Gelen istekleri listele
- `GET /get-sent-requests` - Gönderilen istekleri listele
- `GET /get-friends` - Arkadaş listesi
- `PUT /accept-friend-request/:id` - İsteği kabul et
- `DELETE /reject-friend-request/:id` - İsteği reddet
- `DELETE /cancel-friend-request/:id` - Gönderilen isteği iptal et
- `DELETE /remove-friend/:id` - Arkadaşı kaldır
- `GET /get-pending-requests-count` - Bekleyen istek sayısı

### Yorumlar
- `GET /comments/:poemId` - Şiir yorumlarını listele
- `POST /comment` - Yeni yorum ekle
- `DELETE /comment/:id` - Yorum sil

### Beğeniler ve Bookmarklar
- `POST /liked-poem` - Şiiri beğen
- `DELETE /liked-poem/:id` - Beğeniyi kaldır
- `GET /liked-poems` - Beğenilen şiirleri listele
- `POST /bookmark-poem` - Şiiri bookmark'la
- `DELETE /bookmark-poem/:id` - Bookmark'ı kaldır
- `GET /bookmarked-poems` - Bookmark'lanan şiirleri listele

### WebSocket
- `GET /ws` - WebSocket bağlantısı (real-time notifications)

### Admin
- `GET /get-logs` - Sistem loglarını görüntüle (admin)
- `DELETE /delete-log/:id` - Log kaydını sil (admin)

## Kullanılan Composables (Frontend)

- `useBookActions` - Kitap CRUD işlemleri
- `useConfirmDialog` - Onay dialog yönetimi
- `useErrorHandler` - Hata yönetimi
- `useFriendActions` - Arkadaşlık işlemleri
- `useFriendsList` - Arkadaş listesi yönetimi
- `useLazyLoad` - Lazy loading fonksiyonalitesi
- `useNotification` - Bildirim sistemi
- `usePoemActions` - Şiir CRUD işlemleri
- `useSanitizer` - XSS koruması için input temizleme
- `useSearch` - Arama fonksiyonalitesi
- `useWebSocketFriends` - WebSocket arkadaşlık bildirimleri

## Scripts

### Development
```bash
# Frontend development server
cd frontend && npm run dev

# Backend development server
cd backend && go run main.go
```

### Build
```bash
# Frontend build
cd frontend && npm run build

# Backend build
cd backend && go build -o bin/main main.go
```

### Linting
```bash
# Frontend linting
cd frontend && npm run lint
```

### Deployment Scripts
```bash
# Docker deployment
./scripts/deploy.sh

# Backend only deployment
./scripts/deploy-backend-only.sh

# VPS deployment
./scripts/deploy-vps.sh

# Database backup
./scripts/backup-database.sh

# Database restore
./scripts/restore-database.sh

# Docker cleanup
./scripts/cleanup-docker.sh
```

### Version Release
```bash
# Major version (X.0.0)
./release.sh major backend
./release.sh major frontend

# Minor version (0.X.0)
./release.sh minor backend
./release.sh minor frontend

# Patch version (0.0.X)
./release.sh patch backend
./release.sh patch frontend
```

## Environment Variables

### Backend (.env)
```env
# Database
DB_HOST=localhost
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=poem_blog
DB_PORT=5432
DB_SSLMODE=disable
DB_TIMEZONE=Europe/Istanbul

# Application
APP_PORT=8080
CORS_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
COOKIE_SECURE=false

# Admin
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_admin_password
ADMIN_ROLEID=1
```

### Frontend (Vite)
Frontend ortam değişkenleri Docker build sırasında build argument olarak enjekte edilir:

```env
VITE_API_BASE_URL=http://localhost:8080
```

## Güvenlik Önlemleri

### Backend Güvenlik
- JWT token ile kimlik doğrulama
- Bcrypt ile şifre hashleme
- Rate limiting (global: 100 req/min, auth: 5 req/15min)
- Security headers (XSS, CSRF, clickjacking koruması)
- Input sanitization ve validation
- SQL injection koruması (GORM parametrized queries)
- File upload validation (tip, boyut kontrolü)
- CORS konfigürasyonu
- Request body size limit (10MB)
- Cookie security flags

### Frontend Güvenlik
- DOMPurify ile XSS koruması
- VeeValidate ile form validation
- Input sanitization composable
- HTTPS kullanımı (production)
- Secure cookie flags

## Database Modelleri

### Core Models
- `Admin` - Admin kullanıcıları
- `Author` - Yazarlar
- `Book` - Kitaplar
- `Poem` - Şiirler
- `Comment` - Yorumlar
- `Friendship` - Arkadaşlık ilişkileri
- `Homepage` - Anasayfa içerikleri
- `MihrimahCard` - Mihrimah kartları
- `Reminder` - Hatırlatıcılar
- `Log` - Sistem logları

### İlişkiler
- User → Poems (liked, bookmarked)
- User → Books (read)
- User → Friendships
- Poem → Book → Author
- Poem → Comments

## Production Deployment

### VPS/Cloud Server Deployment

1. **Sunucu Hazırlığı**
```bash
# Docker ve Docker Compose kurulumu
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
sudo usermod -aG docker $USER

# Repository'yi klonlayın
git clone <repository-url>
cd MihrimahSiir-WebSite
```

2. **SSL Sertifikası (Let's Encrypt ile)**
```bash
# Certbot kurulumu
sudo apt-get update
sudo apt-get install certbot

# SSL sertifikası alın
sudo certbot certonly --standalone -d your-domain.com
```

3. **Environment Dosyalarını Yapılandırın**
```bash
cp .env.docker.example .env.docker
cp backend/.env.example backend/.env

# Production değerlerini girin
nano .env.docker
nano backend/.env
```

4. **Deploy Edin**
```bash
chmod +x scripts/deploy-vps.sh
./scripts/deploy-vps.sh
```

5. **Nginx Reverse Proxy (Opsiyonel)**
Eğer SSL ve domain yönetimi için nginx kullanmak isterseniz, ana sunucuda nginx yapılandırın.

### Monitoring ve Bakım
```bash
# Container loglarını görüntüle
docker-compose logs -f

# Container durumu
docker-compose ps

# Veritabanı backup
./scripts/backup-database.sh

# Docker cleanup (disk alanı temizleme)
./scripts/cleanup-docker.sh
```

## Troubleshooting

### Backend başlamıyor
```bash
# Database bağlantısını kontrol edin
docker-compose logs database

# Backend loglarını kontrol edin
docker-compose logs backend

# Environment değişkenlerini kontrol edin
docker-compose exec backend env | grep DB_
```

### Frontend build hatası
```bash
# Node modüllerini temizleyin
cd frontend
rm -rf node_modules package-lock.json
npm install

# Cache'i temizleyin
npm cache clean --force
```

### Database migration sorunları
```bash
# Database'i sıfırlayın (DİKKAT: Tüm veriler silinir!)
docker-compose down -v
docker-compose up -d database
```

### Port çakışması
```bash
# Kullanılan portları kontrol edin
sudo lsof -i :3000
sudo lsof -i :8080
sudo lsof -i :5432

# Çakışan servisi durdurun veya docker-compose.yml'de portları değiştirin
```

## Lisans

Bu proje özel bir projedir. Kullanım hakları hakkında bilgi için proje sahibi ile iletişime geçin.

## İletişim

Proje Sahibi - [GitHub](https://github.com/osmanaydindev)

LinkedIn - [LinkedIn](https://www.linkedin.com/in/osman-ayd%C4%B1n-a22860205/)

Proje Link: [https://github.com/osmanaydin/MihrimahSiir-WebSite](https://github.com/osmanaydindev/MihrimahSiir-WebSite)

Proje Canlı Lİnk: [MihrimahSiir](https://www.mihrimahsiir.com)

## Teşekkürler

- [Fiber](https://gofiber.io/) - Go web framework
- [Vue.js](https://vuejs.org/) - Progressive JavaScript framework
- [Vuetify](https://vuetifyjs.com/) - Material Design component framework
- [GORM](https://gorm.io/) - Go ORM library
- [Vite](https://vitejs.dev/) - Next generation frontend tooling
