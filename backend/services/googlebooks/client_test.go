package googlebooks

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestClient, httptest sunucusuna bağlı istemci kurar. Default()
// sync.Once ile env'den okuduğu için testte doğrudan struct kuruluyor.
func newTestClient(baseURL string) *Client {
	return &Client{
		http:    &http.Client{Timeout: 2 * time.Second},
		baseURL: baseURL,
		sem:     make(chan struct{}, 4),
	}
}

func serve(t *testing.T, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

const turkishVolume = `{
  "totalItems": 1,
  "items": [{
    "id": "abc123",
    "volumeInfo": {
      "title": "Kürk Mantolu Madonna",
      "authors": ["Sabahattin Ali"],
      "publisher": "Yapı Kredi Yayınları",
      "publishedDate": "2016-03-01",
      "description": "<p>Raif Efendi'nin <i>hikâyesi</i>.</p><br>İkinci &quot;paragraf&quot; &amp; devamı.",
      "pageCount": 160,
      "language": "tr",
      "imageLinks": {
        "thumbnail": "http://books.google.com/books/content?id=abc123&zoom=1"
      },
      "industryIdentifiers": [
        {"type": "ISBN_10", "identifier": "9753638029"},
        {"type": "ISBN_13", "identifier": "9789753638029"}
      ]
    }
  }]
}`

func TestFetchByISBNParsesTurkishVolume(t *testing.T) {
	srv := serve(t, turkishVolume)

	meta, err := newTestClient(srv.URL).FetchByISBN(context.Background(), "9789753638029")
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}

	if meta.Title != "Kürk Mantolu Madonna" {
		t.Errorf("başlık = %q", meta.Title)
	}
	if len(meta.Authors) != 1 || meta.Authors[0] != "Sabahattin Ali" {
		t.Errorf("yazarlar = %v", meta.Authors)
	}
	if meta.Publisher != "Yapı Kredi Yayınları" || meta.PublishDate != "2016-03-01" {
		t.Errorf("yayıncı/tarih = %q / %q", meta.Publisher, meta.PublishDate)
	}
	if meta.NumberOfPages != 160 || meta.Language != "tr" {
		t.Errorf("sayfa/dil = %d / %q", meta.NumberOfPages, meta.Language)
	}

	// Açıklama düz metne indirilmeli: etiket kalmamalı, kelimeler
	// yapışmamalı, noktalama öncesi boşluk kalmamalı, entity'ler çözülmeli.
	want := `Raif Efendi'nin hikâyesi. İkinci "paragraf" & devamı.`
	if meta.Description != want {
		t.Errorf("açıklama = %q, beklenen %q", meta.Description, want)
	}

	// Kapak https'e yükseltilmeli; CSP http görseli bloklar.
	wantCover := "https://books.google.com/books/content?id=abc123&zoom=1"
	if meta.CoverURL != wantCover {
		t.Errorf("kapak = %q, beklenen %q", meta.CoverURL, wantCover)
	}
}

// Google bazen sorgulanandan başka bir baskı döndürüyor. Yanlış kitabı
// göstermek, hiç göstermemekten kötü: kullanıcı onu onaylar.
func TestFetchByISBNRejectsDifferentEdition(t *testing.T) {
	srv := serve(t, turkishVolume)

	_, err := newTestClient(srv.URL).FetchByISBN(context.Background(), "9789750718053")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("ErrNotFound beklendi, gelen: %v", err)
	}
}

func TestFetchByISBNEmptyResult(t *testing.T) {
	srv := serve(t, `{"kind":"books#volumes","totalItems":0}`)

	_, err := newTestClient(srv.URL).FetchByISBN(context.Background(), "9789753638029")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("ErrNotFound beklendi, gelen: %v", err)
	}
}

// ISBN_13 taşımayan kayıtlar doğrulanamıyor; eleme yerine kabul edilir,
// yoksa yalnızca ISBN_10 listeleyen kayıtlar boşa düşer.
func TestFetchByISBNAcceptsUnverifiableRecord(t *testing.T) {
	srv := serve(t, `{
      "totalItems": 1,
      "items": [{"id":"x","volumeInfo":{"title":"Eski Baskı","industryIdentifiers":[{"type":"ISBN_10","identifier":"9753638029"}]}}]
    }`)

	meta, err := newTestClient(srv.URL).FetchByISBN(context.Background(), "9789753638029")
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if meta.Title != "Eski Baskı" {
		t.Errorf("başlık = %q", meta.Title)
	}
}

func TestFetchByISBNHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	t.Cleanup(srv.Close)

	_, err := newTestClient(srv.URL).FetchByISBN(context.Background(), "9789753638029")
	if err == nil {
		t.Fatal("kota hatası beklendi")
	}
	// Kota hatası "kayıt yok" demek değil: çağıran yedek kaynağa düşebilsin.
	if errors.Is(err, ErrNotFound) {
		t.Fatalf("429 ErrNotFound'a çevrilmemeli: %v", err)
	}
}
