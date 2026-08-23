// Package openlibrary, openlibrary.org üzerinden ISBN ile kitap
// meta verisi çeker. Repodaki tek dış HTTP entegrasyonlarından biri.
package openlibrary

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// ErrNotFound, ISBN Open Library'de bulunamadığında döner.
var ErrNotFound = errors.New("openlibrary: kayıt bulunamadı")

const (
	defaultBaseURL = "https://openlibrary.org"
	userAgent      = "MihrimahSiir/1.0 (+https://mihrimahsiir.com; osmanacoder@gmail.com)"
	maxBodyBytes   = 512 << 10 // 512 KB
	coverURLPrefix = "https://covers.openlibrary.org/"

	// Toplam bütçe 10s: 5s ana çağrı + 2 x 2.5s açıklama çağrısı.
	primaryTimeout     = 5 * time.Second
	descriptionTimeout = 2500 * time.Millisecond
	totalBudget        = 10 * time.Second
)

// BookMeta, onay ekranında gösterilecek kitap verisi.
type BookMeta struct {
	ISBN          string   `json:"isbn"`
	Title         string   `json:"title"`
	Subtitle      string   `json:"subtitle"`
	Authors       []string `json:"authors"`
	NumberOfPages int      `json:"number_of_pages"`
	CoverURL      string   `json:"cover_url"`
	Publisher     string   `json:"publisher"`
	PublishDate   string   `json:"publish_date"`
	Description   string   `json:"description"`
	EditionKey    string   `json:"edition_key"`
	WorkKey       string   `json:"work_key"`
}

// Client, tek bir http.Client'ı paylaşır. İstek başına client yaratmak
// bağlantı yeniden kullanımını bozar ve yük altında dosya tanıtıcısı sızdırır.
type Client struct {
	http    *http.Client
	baseURL string
	// sem, process genelinde eşzamanlı dış çağrı tavanı. Ani yük altında
	// openlibrary.org'a yüzlerce soket açılmasını engeller.
	sem chan struct{}
}

var (
	defaultClient *Client
	defaultOnce   sync.Once
)

// Default, paket düzeyindeki tek istemciyi döner.
func Default() *Client {
	defaultOnce.Do(func() {
		base := strings.TrimRight(os.Getenv("OPENLIBRARY_BASE_URL"), "/")
		if base == "" {
			base = defaultBaseURL
		}
		defaultClient = &Client{
			http: &http.Client{
				Timeout: 6 * time.Second,
				Transport: &http.Transport{
					MaxIdleConns:          16,
					MaxIdleConnsPerHost:   4,
					IdleConnTimeout:       60 * time.Second,
					ResponseHeaderTimeout: 4 * time.Second,
				},
			},
			baseURL: base,
			sem:     make(chan struct{}, 4),
		}
	})
	return defaultClient
}

// FetchByISBN, normalize edilmiş ISBN-13 ile kitap verisini getirir.
// Açıklama için yapılan ek çağrılar best-effort'tur: başarısız olursa
// Description boş kalır ama fonksiyon yine başarıyla döner.
func (c *Client) FetchByISBN(ctx context.Context, isbn13 string) (*BookMeta, error) {
	ctx, cancel := context.WithTimeout(ctx, totalBudget)
	defer cancel()

	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	meta, err := c.fetchPrimary(ctx, isbn13)
	if err != nil {
		return nil, err
	}

	// Açıklama edition -> work zincirinden geliyor; hata yutulur.
	if desc, workKey := c.fetchDescription(ctx, meta.EditionKey); desc != "" || workKey != "" {
		meta.Description = desc
		meta.WorkKey = workKey
	}

	return meta, nil
}

// primaryResponse, /api/books?jscmd=data yanıtının ihtiyacımız olan kısmı.
type primaryResponse struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	Key           string `json:"key"` // "/books/OL...M"
	NumberOfPages int    `json:"number_of_pages"`
	Authors       []struct {
		Name string `json:"name"`
	} `json:"authors"`
	Publishers []struct {
		Name string `json:"name"`
	} `json:"publishers"`
	PublishDate string `json:"publish_date"`
	Cover       struct {
		Large  string `json:"large"`
		Medium string `json:"medium"`
	} `json:"cover"`
}

func (c *Client) fetchPrimary(ctx context.Context, isbn13 string) (*BookMeta, error) {
	url := fmt.Sprintf("%s/api/books?bibkeys=ISBN:%s&format=json&jscmd=data", c.baseURL, isbn13)

	body, err := c.get(ctx, url, primaryTimeout)
	if err != nil {
		return nil, err
	}

	var payload map[string]primaryResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("openlibrary: yanıt çözümlenemedi: %w", err)
	}

	record, ok := payload["ISBN:"+isbn13]
	if !ok {
		return nil, ErrNotFound
	}

	meta := &BookMeta{
		ISBN:          isbn13,
		Title:         strings.TrimSpace(record.Title),
		Subtitle:      strings.TrimSpace(record.Subtitle),
		NumberOfPages: record.NumberOfPages,
		PublishDate:   strings.TrimSpace(record.PublishDate),
		EditionKey:    record.Key,
	}

	for _, a := range record.Authors {
		if name := strings.TrimSpace(a.Name); name != "" {
			meta.Authors = append(meta.Authors, name)
		}
	}
	if len(record.Publishers) > 0 {
		meta.Publisher = strings.TrimSpace(record.Publishers[0].Name)
	}

	// Kapak yoksa boş bırakılır. /b/isbn/... fallback'i 404 yerine
	// 1x1 boş GIF döndürdüğü için güvenilir değil.
	cover := record.Cover.Large
	if cover == "" {
		cover = record.Cover.Medium
	}
	if strings.HasPrefix(cover, coverURLPrefix) {
		meta.CoverURL = cover
	}

	if meta.Title == "" {
		return nil, ErrNotFound
	}
	return meta, nil
}

// editionResponse, edition kaydından sadece work referansını alır.
type editionResponse struct {
	Works []struct {
		Key string `json:"key"`
	} `json:"works"`
}

// workResponse.Description hem düz string hem {"value": "..."} olabiliyor.
type workResponse struct {
	Description olText `json:"description"`
}

type olText struct {
	Value string
}

func (d *olText) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		d.Value = s
		return nil
	}
	var o struct {
		Value string `json:"value"`
	}
	// Beklenmedik bir şekil gelirse açıklamayı boş bırak, isteği düşürme.
	_ = json.Unmarshal(b, &o)
	d.Value = o.Value
	return nil
}

// fetchDescription, edition -> work zincirini izler. Tüm hatalar yutulur.
func (c *Client) fetchDescription(ctx context.Context, editionKey string) (string, string) {
	if editionKey == "" {
		return "", ""
	}

	editionBody, err := c.get(ctx, c.baseURL+editionKey+".json", descriptionTimeout)
	if err != nil {
		log.Printf("[openlibrary] edition getirilemedi (%s): %v", editionKey, err)
		return "", ""
	}

	var edition editionResponse
	if err := json.Unmarshal(editionBody, &edition); err != nil || len(edition.Works) == 0 {
		return "", ""
	}
	workKey := edition.Works[0].Key

	workBody, err := c.get(ctx, c.baseURL+workKey+".json", descriptionTimeout)
	if err != nil {
		log.Printf("[openlibrary] work getirilemedi (%s): %v", workKey, err)
		return "", workKey
	}

	var work workResponse
	if err := json.Unmarshal(workBody, &work); err != nil {
		return "", workKey
	}
	return strings.TrimSpace(work.Description.Value), workKey
}

// get, ortak HTTP çağrısı: timeout, UA, boyut sınırı ve bağlantı boşaltma.
func (c *Client) get(ctx context.Context, url string, timeout time.Duration) ([]byte, error) {
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Bağlantının yeniden kullanılabilmesi için gövde boşaltılır.
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("openlibrary: beklenmeyen durum kodu %d", resp.StatusCode)
	}

	return io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
}
