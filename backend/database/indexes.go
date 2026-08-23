package database

import "log"

// EnsureIndexes, GORM tag'leriyle ifade edilemeyen Postgres indekslerini
// kurar. Bilerek panic etmez: indeks kurulamazsa uygulama yine ayağa
// kalkar (uygulama katmanında da aynı kontroller var), sadece loglanır.
func EnsureIndexes() {
	statements := []struct {
		name string
		sql  string
	}{
		{
			// Partial unique index: aynı ISBN için aynı anda yalnızca bir
			// "pending" talep olabilir, ama reddedilmiş/onaylanmış talepler
			// birikebilir. (isbn, status) bileşik unique bunu yapamazdı.
			name: "uniq_book_requests_pending_isbn",
			sql: `CREATE UNIQUE INDEX IF NOT EXISTS uniq_book_requests_pending_isbn
			      ON book_requests (isbn) WHERE status = 'pending'`,
		},
		{
			name: "idx_book_requests_user_created",
			sql: `CREATE INDEX IF NOT EXISTS idx_book_requests_user_created
			      ON book_requests (user_id, created_at DESC)`,
		},
	}

	for _, stmt := range statements {
		if err := DB.Exec(stmt.sql).Error; err != nil {
			log.Printf("[database] %s indeksi oluşturulamadı: %v", stmt.name, err)
			continue
		}
	}
}
