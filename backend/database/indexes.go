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
		{
			// Akışın ana sorgusu: görünür yazarların yorumları, tarihe göre.
			name: "idx_comments_admin_created",
			sql: `CREATE INDEX IF NOT EXISTS idx_comments_admin_created
			      ON comments (admin_id, created_at DESC)`,
		},
		{
			name: "idx_admin_liked_poems_admin_created",
			sql: `CREATE INDEX IF NOT EXISTS idx_admin_liked_poems_admin_created
			      ON admin_liked_poems (admin_id, created_at DESC)`,
		},
		{
			name: "idx_user_books_read_admin_created",
			sql: `CREATE INDEX IF NOT EXISTS idx_user_books_read_admin_created
			      ON user_books_read (admin_id, created_at DESC)`,
		},
		// ─── Beğeni/kaydetme join tablolarında mükerrer kayıt ───────────
		// Bu tablolar GORM'un ürettiği many2many tablolarıydı: hiçbir
		// benzersizlik kısıtı yoktu ve Association().Append() çift kayıt
		// ekleyebiliyordu, bu da like_count'u şişiriyordu.
		//
		// AutoMigrate mevcut bir tabloya birincil anahtar EKLEMEZ, o yüzden
		// kısıt burada elle kuruluyor. Önce mükerrerler temizleniyor (ctid
		// ile en eski satır korunur), sonra unique index. Sıra önemli:
		// temizlik yapılmadan index oluşturulamaz.
		{
			name: "dedup_admin_liked_poems",
			sql: `DELETE FROM admin_liked_poems a USING admin_liked_poems b
			      WHERE a.ctid < b.ctid AND a.admin_id = b.admin_id AND a.poem_id = b.poem_id`,
		},
		{
			name: "uniq_admin_liked_poems",
			sql: `CREATE UNIQUE INDEX IF NOT EXISTS uniq_admin_liked_poems
			      ON admin_liked_poems (admin_id, poem_id)`,
		},
		{
			name: "dedup_admin_bookmark_poems",
			sql: `DELETE FROM admin_bookmark_poems a USING admin_bookmark_poems b
			      WHERE a.ctid < b.ctid AND a.admin_id = b.admin_id AND a.poem_id = b.poem_id`,
		},
		{
			name: "uniq_admin_bookmark_poems",
			sql: `CREATE UNIQUE INDEX IF NOT EXISTS uniq_admin_bookmark_poems
			      ON admin_bookmark_poems (admin_id, poem_id)`,
		},
		{
			name: "dedup_user_books_read",
			sql: `DELETE FROM user_books_read a USING user_books_read b
			      WHERE a.ctid < b.ctid AND a.admin_id = b.admin_id AND a.book_id = b.book_id`,
		},
		{
			name: "uniq_user_books_read",
			sql: `CREATE UNIQUE INDEX IF NOT EXISTS uniq_user_books_read
			      ON user_books_read (admin_id, book_id)`,
		},
		{
			// Comment'e created_at sonradan eklendi, mevcut satırlar NULL.
			// Hepsine tek bir geçmiş damga veriliyor: akış
			// (created_at DESC, id DESC) sıraladığı için bunlar yeni
			// içeriğin altında, kendi aralarında id sırasıyla dizilir.
			// NOW() verilseydi bütün eski yorumlar akışın tepesine çıkardı.
			//
			// WHERE koşulu idempotent kılıyor: ikinci çalıştırmada
			// eşleşen satır kalmaz.
			name: "backfill_comments_created_at",
			sql: `UPDATE comments
			      SET created_at = TIMESTAMP '2020-01-01 00:00:00'
			      WHERE created_at IS NULL`,
		},
	}

	for _, stmt := range statements {
		if err := DB.Exec(stmt.sql).Error; err != nil {
			log.Printf("[database] %s indeksi oluşturulamadı: %v", stmt.name, err)
			continue
		}
	}
}
