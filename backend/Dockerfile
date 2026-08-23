# Multi-stage build for Go backend
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Tidy go modules and build
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# ÖNEMLİ: Uygulama dosyaları "./uploads" gibi göreli yollarla kaydeder.
# Bu yüzden WORKDIR, docker-compose'un uploads volume'unu bağladığı
# /app dizini ile aynı olmak zorunda. Aksi halde yüklenen profil resimleri
# container silindiğinde kaybolur.
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Uploads dizini (volume mount edilmezse de var olsun)
RUN mkdir -p /app/uploads/profiles

# .env file will be provided by docker-compose environment variables

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
