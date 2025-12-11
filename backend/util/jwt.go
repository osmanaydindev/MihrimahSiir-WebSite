package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var pkey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
var pkey1 *ecdsa.PublicKey = &pkey.PublicKey

func SetToken(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 gun gecerli
	})

	token, err := claims.SignedString(pkey)
	return token, err
}

type Claims struct {
	jwt.StandardClaims
}

func GetUserWithToken(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return pkey1, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	if err != nil {
		return "token olusmadı", err
	}
	claims := token.Claims.(*Claims)
	return claims.Issuer, err
}
