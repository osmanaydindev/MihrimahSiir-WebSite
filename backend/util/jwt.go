package util

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwtSecret, JWT_SECRET ortam değişkeninden gelir. Anahtar process içinde
// üretilmez; aksi halde her yeniden başlatmada tüm oturumlar düşer.
var jwtSecret []byte

const minSecretLength = 32

// InitJWT, imzalama anahtarını ortam değişkeninden yükler.
// Uygulama başlarken (env yüklendikten sonra) bir kez çağrılmalıdır.
func InitJWT() error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return errors.New("JWT_SECRET tanımlı değil (öneri: openssl rand -hex 32)")
	}
	if len(secret) < minSecretLength {
		return fmt.Errorf("JWT_SECRET en az %d karakter olmalı", minSecretLength)
	}
	jwtSecret = []byte(secret)
	return nil
}

func SetToken(issuer string) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("jwt: imzalama anahtarı yüklenmedi")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 gun gecerli
	})

	token, err := claims.SignedString(jwtSecret)
	return token, err
}

type Claims struct {
	jwt.StandardClaims
}

func GetUserWithToken(cookie string) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("jwt: imzalama anahtarı yüklenmedi")
	}

	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Algoritma karıştırma saldırılarına karşı imza yöntemini doğrula
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("beklenmeyen imza yöntemi: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("token gecersiz")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", errors.New("token claim'leri okunamadi")
	}
	return claims.Issuer, nil
}
