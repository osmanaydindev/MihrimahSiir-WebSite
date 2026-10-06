// Package bookmeta, ISBN -> kitap meta verisi aramasını tek kapıya
// indirir. Controller'lar hangi dış servisin cevap verdiğini bilmek
// zorunda kalmasın diye var.
//
// Sıra bilinçli: Google Books birincil, Open Library yedek. Open Library'nin
// Türkçe baskı kapsaması zayıf (çoğu Türkçe ISBN'de boş cevap dönüyor),
// Google Books'ta ise aynı kitaplar büyük oranda bulunuyor. Ama Google
// anahtar gerektiriyor ve kotası var; anahtar yoksa veya kota dolmuşsa
// Open Library hâlâ bir şey bulabilir, o yüzden ikisi birden duruyor.
package bookmeta

import (
	"context"
	"errors"
	"log"

	"backend/services/googlebooks"
	"backend/services/openlibrary"
)

// ErrNotFound, ISBN hiçbir kaynakta bulunamadığında döner.
var ErrNotFound = errors.New("bookmeta: kayıt bulunamadı")

// Kaynak adları. Talep kaydında ve önizlemede saklanıyor; "bu bilgi
// nereden geldi" sorusunun cevabı admin için önemli, çünkü iki servisin
// güvenilirliği farklı.
const (
	SourceGoogleBooks = "google_books"
	SourceOpenLibrary = "open_library"
)

// Meta, kaynaktan bağımsız kitap anlık görüntüsü.
type Meta struct {
	ISBN          string
	Title         string
	Subtitle      string
	Authors       []string
	NumberOfPages int
	CoverURL      string
	Publisher     string
	PublishDate   string
	Description   string
	// Source, veriyi döndüren servis (Source* sabitlerinden biri).
	Source string
	// OpenLibraryKey yalnızca Open Library cevabında dolu olur.
	OpenLibraryKey string
}

// Lookup, ISBN-13 için meta veriyi getirir.
//
// Kaynaklardan biri "bulunamadı" derse diğeri denenir; ağ/kota hatası da
// aynı şekilde yedeğe düşer (loglanır, isteği düşürmez). İkisi de boş
// dönerse ErrNotFound. Dönen hata ErrNotFound değilse "şu an ulaşılamıyor"
// demektir — çağıran bunu kullanıcıya farklı mesajla gösterebilir.
func Lookup(ctx context.Context, isbn13 string) (*Meta, error) {
	var transportErr error

	if meta, err := googlebooks.Default().FetchByISBN(ctx, isbn13); err == nil {
		return fromGoogle(meta), nil
	} else if !errors.Is(err, googlebooks.ErrNotFound) {
		transportErr = err
		log.Printf("[bookmeta] Google Books çağrısı başarısız (%s): %v", isbn13, err)
	}

	if meta, err := openlibrary.Default().FetchByISBN(ctx, isbn13); err == nil {
		return fromOpenLibrary(meta), nil
	} else if !errors.Is(err, openlibrary.ErrNotFound) {
		transportErr = err
		log.Printf("[bookmeta] Open Library çağrısı başarısız (%s): %v", isbn13, err)
	}

	// Kayıt gerçekten yok mu, yoksa iki servise de ulaşamadık mı? İkisi
	// kullanıcıya farklı şey söylemeyi gerektiriyor.
	if transportErr != nil {
		return nil, transportErr
	}
	return nil, ErrNotFound
}

func fromGoogle(meta *googlebooks.BookMeta) *Meta {
	return &Meta{
		ISBN:          meta.ISBN,
		Title:         meta.Title,
		Subtitle:      meta.Subtitle,
		Authors:       meta.Authors,
		NumberOfPages: meta.NumberOfPages,
		CoverURL:      meta.CoverURL,
		Publisher:     meta.Publisher,
		PublishDate:   meta.PublishDate,
		Description:   meta.Description,
		Source:        SourceGoogleBooks,
	}
}

func fromOpenLibrary(meta *openlibrary.BookMeta) *Meta {
	return &Meta{
		ISBN:           meta.ISBN,
		Title:          meta.Title,
		Subtitle:       meta.Subtitle,
		Authors:        meta.Authors,
		NumberOfPages:  meta.NumberOfPages,
		CoverURL:       meta.CoverURL,
		Publisher:      meta.Publisher,
		PublishDate:    meta.PublishDate,
		Description:    meta.Description,
		Source:         SourceOpenLibrary,
		OpenLibraryKey: meta.EditionKey,
	}
}
