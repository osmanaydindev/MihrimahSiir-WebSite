// Package googlebooks, Google Books API üzerinden ISBN ile kitap meta
// verisi çeker. Open Library'nin Türkçe baskı kapsaması zayıf olduğu için
// birincil kaynak burasıdır; openlibrary paketi yedekte kalır.
//
// Disiplin openlibrary paketiyle aynı: paket düzeyinde tek http.Client,
// context timeout'ları, io.LimitReader ile gövde sınırı, gövdenin kapatma
// öncesi boşaltılması ve process genelinde eşzamanlı çağrı tavanı.
package googlebooks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

// ErrNotFound, ISBN Google Books'ta bulunamadığında döner.
var ErrNotFound = errors.New("googlebooks: kayıt bulunamadı")

const (
	defaultBaseURL = "https://www.googleapis.com/books/v1"
	userAgent      = "MihrimahSiir/1.0 (+https://mihrimahsiir.com; osmanacoder@gmail.com)"
	maxBodyBytes   = 512 << 10 // 512 KB
	requestTimeout = 5 * time.Second
	totalBudget    = 6 * time.Second
)

// BookMeta, onay/önizleme ekranında gösterilecek kitap verisi.
// Alan adları openlibrary.BookMeta ile kasıtlı olarak aynı; iki kaynak
// bookmeta paketinde tek tipe dönüştürülüyor.
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
	VolumeID      string   `json:"volume_id"`
	Language      string   `json:"language"`
}

type Client struct {
	http    *http.Client
	baseURL string
	apiKey  string
	// sem, process genelinde eşzamanlı dış çağrı tavanı.
	sem chan struct{}
}

var (
	defaultClient *Client
	defaultOnce   sync.Once
)

// Default, paket düzeyindeki tek istemciyi döner.
func Default() *Client {
	defaultOnce.Do(func() {
		base := strings.TrimRight(os.Getenv("GOOGLE_BOOKS_BASE_URL"), "/")
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
			apiKey:  strings.TrimSpace(os.Getenv("GOOGLE_BOOKS_API_KEY")),
			sem:     make(chan struct{}, 4),
		}
	})
	return defaultClient
}

// HasAPIKey, anahtarın tanımlı olup olmadığını söyler. Anahtarsız istekler
// de çalışıyor ama paylaşımlı anonim kotaya düştüğü için sık 429 alır;
// çağıran bu durumda yedek kaynağa geçmeyi tercih edebilir.
func (c *Client) HasAPIKey() bool {
	return c.apiKey != ""
}

// imageLinks, Google'ın döndüğü kapak boyutları (büyükten küçüğe).
type imageLinks struct {
	ExtraLarge     string `json:"extraLarge"`
	Large          string `json:"large"`
	Medium         string `json:"medium"`
	Small          string `json:"small"`
	Thumbnail      string `json:"thumbnail"`
	SmallThumbnail string `json:"smallThumbnail"`
}

// industryIdentifier, kayıttaki ISBN_10 / ISBN_13 / OTHER girdileri.
type industryIdentifier struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

// volumesResponse, /volumes yanıtının ihtiyacımız olan kısmı.
type volumesResponse struct {
	TotalItems int `json:"totalItems"`
	Items      []struct {
		ID         string `json:"id"`
		VolumeInfo struct {
			Title               string               `json:"title"`
			Subtitle            string               `json:"subtitle"`
			Authors             []string             `json:"authors"`
			Publisher           string               `json:"publisher"`
			PublishedDate       string               `json:"publishedDate"`
			Description         string               `json:"description"`
			PageCount           int                  `json:"pageCount"`
			Language            string               `json:"language"`
			ImageLinks          imageLinks           `json:"imageLinks"`
			IndustryIdentifiers []industryIdentifier `json:"industryIdentifiers"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

// FetchByISBN, normalize edilmiş ISBN-13 ile kitap verisini getirir.
func (c *Client) FetchByISBN(ctx context.Context, isbn13 string) (*BookMeta, error) {
	ctx, cancel := context.WithTimeout(ctx, totalBudget)
	defer cancel()

	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// country parametresi zorunlu gibi davranıyor: bazı sunucu
	// lokasyonlarında eksikse API 403 "unsupported country" dönüyor.
	params := url.Values{}
	params.Set("q", "isbn:"+isbn13)
	params.Set("country", "TR")
	params.Set("maxResults", "1")
	if c.apiKey != "" {
		params.Set("key", c.apiKey)
	}

	body, err := c.get(ctx, c.baseURL+"/volumes?"+params.Encode())
	if err != nil {
		return nil, err
	}

	var payload volumesResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("googlebooks: yanıt çözümlenemedi: %w", err)
	}
	if payload.TotalItems == 0 || len(payload.Items) == 0 {
		return nil, ErrNotFound
	}

	info := payload.Items[0].VolumeInfo
	title := strings.TrimSpace(info.Title)
	if title == "" {
		return nil, ErrNotFound
	}

	// Sorgulanan ISBN ile dönen kaydın ISBN'i tutmuyorsa bu başka bir
	// baskıdır. "Makul görünen yanlış kitap" göstermek, hiç göstermemekten
	// kötü: kullanıcı yanlış kitabı onaylar. Doğrulanamıyorsa (kayıtta
	// ISBN_13 yoksa) kabul edilir.
	if ids := isbn13s(info.IndustryIdentifiers); len(ids) > 0 && !slices.Contains(ids, isbn13) {
		return nil, ErrNotFound
	}

	meta := &BookMeta{
		ISBN:          isbn13,
		Title:         title,
		Subtitle:      strings.TrimSpace(info.Subtitle),
		NumberOfPages: info.PageCount,
		Publisher:     strings.TrimSpace(info.Publisher),
		PublishDate:   strings.TrimSpace(info.PublishedDate),
		Description:   stripHTML(info.Description),
		VolumeID:      payload.Items[0].ID,
		Language:      strings.TrimSpace(info.Language),
		CoverURL:      pickCover(info.ImageLinks),
	}

	for _, a := range info.Authors {
		if name := strings.TrimSpace(a); name != "" {
			meta.Authors = append(meta.Authors, name)
		}
	}

	return meta, nil
}

// isbn13s, kayıttaki ISBN-13 değerlerini (tireleri atılmış halde) döner.
func isbn13s(ids []industryIdentifier) []string {
	var out []string
	for _, id := range ids {
		if id.Type != "ISBN_13" {
			continue
		}
		clean := strings.NewReplacer("-", "", " ", "").Replace(id.Identifier)
		if clean != "" {
			out = append(out, clean)
		}
	}
	return out
}

// pickCover, en büyük kapağı seçer ve kullanılabilir hâle getirir.
//
// İki düzeltme şart:
//   - http -> https: sitenin CSP'si img-src için 'self' data: https:
//     tanımlı, Google ise thumbnail'leri http olarak dönüyor.
//   - zoom/edge: arama sonucunda yalnızca thumbnail (zoom=1) geliyor ve o
//     128x198 piksel — kart için bile bulanık. Aynı uç nokta zoom=2 ile
//     300x464 (~20 KB) veriyor. edge=curl ise kapağa sahte "kıvrık sayfa"
//     efekti bindiriyor, atılıyor.
func pickCover(links imageLinks) string {
	for _, candidate := range []string{
		links.ExtraLarge, links.Large, links.Medium,
		links.Small, links.Thumbnail, links.SmallThumbnail,
	} {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		if strings.HasPrefix(candidate, "http://") {
			candidate = "https://" + strings.TrimPrefix(candidate, "http://")
		}
		if !strings.HasPrefix(candidate, "https://") {
			continue
		}
		return upgradeCover(candidate)
	}
	return ""
}

// upgradeCover, books.google.com/books/content bağlantısını büyütür.
// Tanımadığı bir URL'e dokunmaz.
func upgradeCover(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	q := parsed.Query()
	if q.Get("zoom") == "" && !q.Has("edge") {
		return raw
	}
	if zoom := q.Get("zoom"); zoom == "" || zoom == "1" {
		q.Set("zoom", "2")
	}
	q.Del("edge")
	parsed.RawQuery = q.Encode()
	return parsed.String()
}

// stripHTML, Google açıklamalarındaki <p>/<br>/<i> gibi etiketleri atar.
// Açıklama şablonlarda ve admin panelinde düz metin olarak gösteriliyor;
// etiketli metin oraya hiç ulaşmasın.
func stripHTML(input string) string {
	if input == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(input))
	depth := 0
	for _, r := range input {
		switch {
		case r == '<':
			depth++
		case r == '>' && depth > 0:
			depth--
			// Etiket sınırı kelimeleri birbirine yapıştırmasın.
			b.WriteRune(' ')
		case depth == 0:
			b.WriteRune(r)
		}
	}
	// Etiket temizliğinden kalan çoklu boşlukları sadeleştir.
	out := strings.Join(strings.Fields(b.String()), " ")

	// Etiket sınırına koyduğumuz boşluk noktalamanın önüne düşmüş olabilir
	// ("hikâyesi ." gibi); bunları geri topla.
	for _, p := range []string{".", ",", "!", "?", ":", ";", ")", "…", "”", "\""} {
		out = strings.ReplaceAll(out, " "+p, p)
	}

	// Google açıklamaları HTML entity içeriyor (&quot;, &amp;, &#39;).
	// Etiketler atıldıktan sonra çözülüyor; önce çözülse "&lt;b&gt;"
	// gerçek etikete dönüşüp bir sonraki adımı atlatırdı.
	return strings.TrimSpace(html.UnescapeString(out))
}

// get, ortak HTTP çağrısı: timeout, UA, boyut sınırı ve bağlantı boşaltma.
func (c *Client) get(ctx context.Context, requestURL string) ([]byte, error) {
	reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, requestURL, nil)
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
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("googlebooks: beklenmeyen durum kodu %d", resp.StatusCode)
	}

	return io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
}
