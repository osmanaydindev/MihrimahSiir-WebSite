package security

import "testing"

func TestNormalizeISBN(t *testing.T) {
	v := NewValidator()

	valid := []struct {
		in   string
		want string
	}{
		// Kürk Mantolu Madonna — ISBN-13 ve tireli hali aynı sonucu vermeli
		{"9789753638029", "9789753638029"},
		{"978-975-363-802-9", "9789753638029"},
		{"978 975 363 802 9", "9789753638029"},
		// 1984 (Penguin) ISBN-10 -> ISBN-13
		{"0140328726", "9780140328721"},
		{"9780140328721", "9780140328721"},
		// Son hanesi X olan ISBN-10
		{"080442957X", "9780804429573"},
	}
	for _, tc := range valid {
		got, err := v.NormalizeISBN(tc.in)
		if err != nil {
			t.Errorf("NormalizeISBN(%q) beklenmedik hata: %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("NormalizeISBN(%q) = %q, beklenen %q", tc.in, got, tc.want)
		}
	}

	invalid := []string{
		"",
		"1234",
		"9789753638028",  // bozuk ISBN-13 checksum
		"0140328727",     // bozuk ISBN-10 checksum
		"1234567890123",  // 978/979 ile başlamıyor
		"97897536380299", // 14 hane
		"978975363802X",  // ISBN-13'te X olamaz
	}
	for _, in := range invalid {
		if got, err := v.NormalizeISBN(in); err == nil {
			t.Errorf("NormalizeISBN(%q) hata vermeliydi, %q döndü", in, got)
		}
	}
}

func TestValidateBookCommunity(t *testing.T) {
	v := NewValidator()
	for _, ok := range []int{1, 2, 3} {
		if err := v.ValidateBookCommunity(ok); err != nil {
			t.Errorf("community %d geçerli olmalı: %v", ok, err)
		}
	}
	for _, bad := range []int{0, 4, -1} {
		if err := v.ValidateBookCommunity(bad); err == nil {
			t.Errorf("community %d reddedilmeliydi", bad)
		}
	}
	// Şiirler hâlâ sadece 1|2 kabul etmeli (community=3 şiirde görünmez olurdu)
	if err := v.ValidateCommunity(3); err == nil {
		t.Error("ValidateCommunity(3) reddedilmeliydi")
	}
}

func TestSanitizeStringTruncatesOnRuneBoundary(t *testing.T) {
	s := NewSanitizer()
	// Her biri 2 byte olan Türkçe karakterler
	got := s.SanitizeString("şşşşş", 3)
	if got != "şşş" {
		t.Errorf("beklenen %q, gelen %q", "şşş", got)
	}
	for i, r := range got {
		if r == '�' {
			t.Errorf("bozuk rune %d. konumda: %q", i, got)
		}
	}
}
