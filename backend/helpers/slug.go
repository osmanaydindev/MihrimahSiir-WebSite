package helpers

import (
	"backend/database"
	"fmt"
	"strings"
)

// turkishSlugMap, controllers içinde 4 ayrı yere kopyalanmış olan
// karakter tablosunun tek kopyası.
var turkishSlugMap = map[rune]rune{
	' ':  '-',
	'ç':  'c',
	'ğ':  'g',
	'ı':  'i',
	'ö':  'o',
	'ş':  's',
	'ü':  'u',
	'Ç':  'C',
	'Ğ':  'G',
	'İ':  'I',
	'Ö':  'O',
	'Ş':  'S',
	'Ü':  'U',
	'â':  'a',
	'Â':  'A',
	'\'': '-',
}

// Slugify, Türkçe karakterleri sadeleştirip küçük harfli slug üretir.
func Slugify(s string) string {
	var b strings.Builder
	for _, char := range s {
		if replacement, ok := turkishSlugMap[char]; ok {
			b.WriteRune(replacement)
		} else {
			b.WriteRune(char)
		}
	}
	return strings.ToLower(b.String())
}

// UniqueBookSlug, books tablosunda çakışmayan bir slug döner.
// books.slug üzerinde unique kısıt yok ve GetBook slug'ı First() ile
// arıyor; aynı isimli ikinci kitap tekilleştirilmezse slug üzerinden
// asla erişilemez hale gelir.
func UniqueBookSlug(name string, excludeBookID uint) string {
	base := Slugify(name)
	if base == "" {
		base = "kitap"
	}

	candidate := base
	for i := 2; i < 1000; i++ {
		var count int64
		query := database.DB.Table("books").Where("slug = ?", candidate)
		if excludeBookID > 0 {
			query = query.Where("id <> ?", excludeBookID)
		}
		query.Count(&count)

		if count == 0 {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return candidate
}
