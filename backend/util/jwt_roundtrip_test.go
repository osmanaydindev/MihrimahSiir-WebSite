package util

import (
	"os"
	"testing"
)

// Aynı JWT_SECRET ile "yeniden başlatma" simüle edilir: eski token geçerli kalmalı.
func TestTokenSurvivesRestart(t *testing.T) {
	os.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	if err := InitJWT(); err != nil {
		t.Fatal(err)
	}
	token, err := SetToken("42")
	if err != nil {
		t.Fatal(err)
	}

	// restart
	jwtSecret = nil
	if err := InitJWT(); err != nil {
		t.Fatal(err)
	}
	id, err := GetUserWithToken(token)
	if err != nil {
		t.Fatalf("restart sonrası token gecersiz: %v", err)
	}
	if id != "42" {
		t.Fatalf("beklenen 42, gelen %q", id)
	}

	// farklı secret ile doğrulanmamalı
	os.Setenv("JWT_SECRET", "ffffffffffffffffffffffffffffffff")
	if err := InitJWT(); err != nil {
		t.Fatal(err)
	}
	if _, err := GetUserWithToken(token); err == nil {
		t.Fatal("farklı secret ile token kabul edildi")
	}
}

func TestInitJWTRejectsMissingOrShortSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "")
	if err := InitJWT(); err == nil {
		t.Fatal("boş JWT_SECRET kabul edildi")
	}
	os.Setenv("JWT_SECRET", "kisa")
	if err := InitJWT(); err == nil {
		t.Fatal("kısa JWT_SECRET kabul edildi")
	}
}
