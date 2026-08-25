package controllers

import (
	"time"

	"backend/database"
	"backend/helpers"
	"backend/models"

	"github.com/gofiber/fiber/v2"
)

// ─── Akış ───────────────────────────────────────────────────────────────────
//
// Görünürlük iki ekseni BİRDEN karşılamak zorunda:
//
//  1. Yazar ekseni — içeriği yazan kişiyi görebiliyor muyum?
//     arkadaşım VEYA profili herkese açık.
//  2. İçerik ekseni — içeriğin asıldığı kitabı/şiiri görebiliyor muyum?
//     (community filtreleri)
//
// İkincisi atlanırsa özel kitap adları sızar. Repoda emsali var:
// GetReadsBooksPaginated bu filtreyi hiç uygulamıyor.
//
// Akış bilerek "okundu" takibi yapmıyor ve zaman penceresi uygulamıyor:
// yeni içerik kalmadığında geçmişe doğru sayfalamaya devam eder, böylece
// akış hiçbir zaman boş görünmez.

type feedRow struct {
	Kind      string    `json:"kind"` // comment | poem_like | book_read
	RefID     uint      `json:"ref_id"`
	AdminID   uint      `json:"admin_id"`
	CreatedAt time.Time `json:"created_at"`
}

type feedAuthor struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	ProfileImage string `json:"profile_image"`
}

// FeedItem, istemciye giden birleşik akış öğesi. Yalnızca `kind`e ait alanlar
// doldurulur, geri kalanı boş bırakılır.
type FeedItem struct {
	Kind      string     `json:"kind"`
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	Author    feedAuthor `json:"author"`

	// kind == comment
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	Page       *int   `json:"page,omitempty"`
	BookID     uint   `json:"book_id,omitempty"`
	BookName   string `json:"book_name,omitempty"`
	BookSlug   string `json:"book_slug,omitempty"`
	LikeCount  int    `json:"like_count"`
	IsLiked    bool   `json:"is_liked"`
	IsSaved    bool   `json:"is_saved"`

	// kind == poem_like
	PoemTitle  string `json:"poem_title,omitempty"`
	PoemSlug   string `json:"poem_slug,omitempty"`
	PoemAuthor string `json:"poem_author,omitempty"`

	// kind == book_read (BookName/BookSlug yeniden kullanılıyor)
}

// visibleAuthorIDs, izleyicinin içeriğini görebileceği kullanıcı id'lerini
// döner: arkadaşlar (kendisi dahil) + profili herkese açık olanlar.
func visibleAuthorIDs(viewerID uint) []uint {
	seen := map[uint]bool{}
	ids := []uint{}
	for _, id := range helpers.FriendIDs(viewerID) {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	for _, id := range helpers.PublicUserIDs() {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids
}

// GetFeed — GET /feed?offset=&limit=
func GetFeed(c *fiber.Ctx) error {
	viewerID := GetUserId(c)
	if viewerID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticated"})
	}
	roleID, _ := helpers.GetUserRole(c)
	p := helpers.GetPaginationParams(c)

	friendIDs := helpers.FriendIDs(viewerID)
	authorIDs := visibleAuthorIDs(viewerID)
	if len(authorIDs) == 0 {
		return c.JSON(helpers.CreatePaginationResponse([]FeedItem{}, 0, p.Offset, p.Limit))
	}

	// Yorumlar: arkadaşınsa her zaman görünür; arkadaş değilse yalnızca
	// profili public OLAN ve YENİ kural altında yazılmış olanlar.
	commentWhere := database.DB.Model(&models.Comment{}).
		Select("'comment' AS kind, comments.id AS ref_id, comments.admin_id, comments.created_at").
		Joins("JOIN admins ON admins.id = comments.admin_id").
		Joins("JOIN books ON books.id = comments.book_id").
		Where("comments.is_deleted = ?", false).
		Where("comments.admin_id IN ?", authorIDs).
		Where("(comments.admin_id IN ? OR (admins.is_private = ? AND comments.allow_public_feed = ?))",
			friendIDs, false, true)
	commentWhere = helpers.ApplyBookCommunityFilter(commentWhere, roleID, viewerID)

	// Aktivite: zaman damgası olmayan (migration öncesi) satırlar akışa
	// girmez — ne zaman olduklarını bilmiyoruz.
	poemLikeWhere := database.DB.Table("admin_liked_poems").
		Select("'poem_like' AS kind, admin_liked_poems.poem_id AS ref_id, admin_liked_poems.admin_id, admin_liked_poems.created_at").
		Joins("JOIN poems ON poems.id = admin_liked_poems.poem_id").
		Where("admin_liked_poems.created_at IS NOT NULL").
		Where("admin_liked_poems.admin_id IN ?", authorIDs).
		Where("poems.is_deleted = ?", false)
	if roleID != 1 && roleID != 2 {
		poemLikeWhere = poemLikeWhere.Where("poems.community = ?", 2)
	}

	bookReadWhere := database.DB.Table("user_books_read").
		Select("'book_read' AS kind, user_books_read.book_id AS ref_id, user_books_read.admin_id, user_books_read.created_at").
		Joins("JOIN books ON books.id = user_books_read.book_id").
		Where("user_books_read.created_at IS NOT NULL").
		Where("user_books_read.admin_id IN ?", authorIDs)
	bookReadWhere = helpers.ApplyBookCommunityFilter(bookReadWhere, roleID, viewerID)

	union := database.DB.
		Raw("? UNION ALL ? UNION ALL ?", commentWhere, poemLikeWhere, bookReadWhere)

	var total int64
	database.DB.Table("(?) AS feed", union).Count(&total)

	var rows []feedRow
	if err := database.DB.
		Table("(?) AS feed", union).
		Order("created_at DESC, ref_id DESC").
		Offset(p.Offset).
		Limit(p.Limit).
		Scan(&rows).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "feed could not be loaded"})
	}

	items := hydrateFeed(rows, viewerID)
	return c.JSON(helpers.CreatePaginationResponse(items, total, p.Offset, p.Limit))
}

// hydrateFeed, hafif akış satırlarını türlerine göre toplu sorgularla
// zenginleştirir. Satır başına sorgu atılmıyor (N+1 yok).
func hydrateFeed(rows []feedRow, viewerID uint) []FeedItem {
	items := make([]FeedItem, 0, len(rows))
	if len(rows) == 0 {
		return items
	}

	commentIDs, poemIDs, bookIDs, adminIDs := []uint{}, []uint{}, []uint{}, []uint{}
	for _, r := range rows {
		adminIDs = append(adminIDs, r.AdminID)
		switch r.Kind {
		case "comment":
			commentIDs = append(commentIDs, r.RefID)
		case "poem_like":
			poemIDs = append(poemIDs, r.RefID)
		case "book_read":
			bookIDs = append(bookIDs, r.RefID)
		}
	}

	// Yazarlar
	authors := map[uint]feedAuthor{}
	if len(adminIDs) > 0 {
		var admins []models.Admin
		database.DB.Select("id, username, profile_image").Where("id IN ?", adminIDs).Find(&admins)
		for _, a := range admins {
			authors[a.ID] = feedAuthor{ID: a.ID, Username: a.Username, ProfileImage: a.ProfileImage}
		}
	}

	// Yorumlar + bağlı kitap
	comments := map[uint]models.Comment{}
	commentBook := map[uint]models.Book{}
	likeCounts := map[uint]int{}
	liked := map[uint]bool{}
	saved := map[uint]bool{}
	if len(commentIDs) > 0 {
		var cs []models.Comment
		database.DB.Where("id IN ?", commentIDs).Find(&cs)
		bIDs := []uint{}
		for _, cm := range cs {
			comments[cm.ID] = cm
			bIDs = append(bIDs, cm.BookID)
		}
		if len(bIDs) > 0 {
			var bs []models.Book
			database.DB.Select("id, name, slug").Where("id IN ?", bIDs).Find(&bs)
			byID := map[uint]models.Book{}
			for _, b := range bs {
				byID[b.ID] = b
			}
			for _, cm := range cs {
				commentBook[cm.ID] = byID[cm.BookID]
			}
		}

		type countRow struct {
			CommentID uint
			N         int
		}
		var counts []countRow
		database.DB.Table("comment_likes").
			Select("comment_id, COUNT(*) AS n").
			Where("comment_id IN ?", commentIDs).
			Group("comment_id").Scan(&counts)
		for _, cr := range counts {
			likeCounts[cr.CommentID] = cr.N
		}

		myLikes := []uint{}
		database.DB.Table("comment_likes").Where("admin_id = ? AND comment_id IN ?", viewerID, commentIDs).
			Pluck("comment_id", &myLikes)
		for _, id := range myLikes {
			liked[id] = true
		}

		mySaves := []uint{}
		database.DB.Table("comment_saves").Where("admin_id = ? AND comment_id IN ?", viewerID, commentIDs).
			Pluck("comment_id", &mySaves)
		for _, id := range mySaves {
			saved[id] = true
		}
	}

	poems := map[uint]models.Poem{}
	if len(poemIDs) > 0 {
		var ps []models.Poem
		database.DB.Select("id, title, slug, author").Where("id IN ?", poemIDs).Find(&ps)
		for _, p := range ps {
			poems[p.ID] = p
		}
	}

	books := map[uint]models.Book{}
	if len(bookIDs) > 0 {
		var bs []models.Book
		database.DB.Select("id, name, slug").Where("id IN ?", bookIDs).Find(&bs)
		for _, b := range bs {
			books[b.ID] = b
		}
	}

	for _, r := range rows {
		item := FeedItem{Kind: r.Kind, ID: r.RefID, CreatedAt: r.CreatedAt, Author: authors[r.AdminID]}
		switch r.Kind {
		case "comment":
			cm, ok := comments[r.RefID]
			if !ok {
				continue
			}
			bk := commentBook[r.RefID]
			item.Title, item.Content, item.Page = cm.Title, cm.Content, cm.Page
			item.BookID, item.BookName, item.BookSlug = bk.ID, bk.Name, bk.Slug
			item.LikeCount, item.IsLiked, item.IsSaved = likeCounts[r.RefID], liked[r.RefID], saved[r.RefID]
		case "poem_like":
			pm, ok := poems[r.RefID]
			if !ok {
				continue
			}
			item.PoemTitle, item.PoemSlug, item.PoemAuthor = pm.Title, pm.Slug, pm.Author
		case "book_read":
			bk, ok := books[r.RefID]
			if !ok {
				continue
			}
			item.BookID, item.BookName, item.BookSlug = bk.ID, bk.Name, bk.Slug
		}
		items = append(items, item)
	}
	return items
}

// ─── Alıntı beğenme / kaydetme ──────────────────────────────────────────────
//
// Uçlar aktörü çerezden alıyor. Mevcut /add-poem-to-liked/:id deseni bilerek
// kopyalanmadı: oradaki :id eylemi yapanın kendi id'si ve handler
// "GetUserId(c) != id" diye kontrol ediyor — gereksiz ve hataya açık.

func toggleCommentRelation(c *fiber.Ctx, table string, add bool) error {
	viewerID := GetUserId(c)
	if viewerID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticated"})
	}
	commentID, err := c.ParamsInt("comment_id")
	if err != nil || commentID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid comment id"})
	}

	if !add {
		if err := database.DB.Exec(
			"DELETE FROM "+table+" WHERE admin_id = ? AND comment_id = ?", viewerID, commentID,
		).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "could not update"})
		}
		return c.JSON(fiber.Map{"message": "ok"})
	}

	// Beğenmeden önce yorumun izleyiciye gerçekten görünür olduğunu doğrula:
	// aksi halde id deneyerek gizli içeriğin varlığı öğrenilebilir, üstelik
	// görülemeyen içerikle etkileşilebilir.
	//
	// Kontrol akıştaki filtreyle BİREBİR aynı olmak zorunda — iki eksen
	// birden. İlk yazımda yalnız yazar ekseni vardı ve özel bir kitaba
	// yazılmış alıntı akışta gizliyken beğenilebiliyordu.
	roleID, _ := helpers.GetUserRole(c)
	var visible int64
	q := database.DB.Model(&models.Comment{}).
		Joins("JOIN admins ON admins.id = comments.admin_id").
		Joins("JOIN books ON books.id = comments.book_id").
		Where("comments.id = ? AND comments.is_deleted = ?", commentID, false).
		Where("(comments.admin_id IN ? OR (admins.is_private = ? AND comments.allow_public_feed = ?))",
			helpers.FriendIDs(viewerID), false, true)
	q = helpers.ApplyBookCommunityFilter(q, roleID, viewerID)
	q.Count(&visible)
	if visible == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
	}

	// ON CONFLICT: bileşik birincil anahtar sayesinde çift kayıt oluşmaz.
	if err := database.DB.Exec(
		"INSERT INTO "+table+" (admin_id, comment_id, created_at) VALUES (?, ?, ?) ON CONFLICT DO NOTHING",
		viewerID, commentID, time.Now(),
	).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "could not update"})
	}
	return c.JSON(fiber.Map{"message": "ok"})
}

func LikeComment(c *fiber.Ctx) error   { return toggleCommentRelation(c, "comment_likes", true) }
func UnlikeComment(c *fiber.Ctx) error { return toggleCommentRelation(c, "comment_likes", false) }
func SaveComment(c *fiber.Ctx) error   { return toggleCommentRelation(c, "comment_saves", true) }
func UnsaveComment(c *fiber.Ctx) error { return toggleCommentRelation(c, "comment_saves", false) }

// GetSavedComments — kullanıcının kaydettiği alıntılar, kaydetme tarihine göre.
func GetSavedComments(c *fiber.Ctx) error {
	viewerID := GetUserId(c)
	if viewerID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticated"})
	}
	roleID, _ := helpers.GetUserRole(c)
	p := helpers.GetPaginationParams(c)

	base := database.DB.Table("comment_saves").
		Select("'comment' AS kind, comment_saves.comment_id AS ref_id, comments.admin_id, comment_saves.created_at").
		Joins("JOIN comments ON comments.id = comment_saves.comment_id").
		Joins("JOIN books ON books.id = comments.book_id").
		Where("comment_saves.admin_id = ?", viewerID).
		Where("comments.is_deleted = ?", false)
	base = helpers.ApplyBookCommunityFilter(base, roleID, viewerID)

	var total int64
	database.DB.Table("(?) AS saved", base).Count(&total)

	var rows []feedRow
	if err := database.DB.Table("(?) AS saved", base).
		Order("created_at DESC, ref_id DESC").
		Offset(p.Offset).Limit(p.Limit).Scan(&rows).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "could not load"})
	}

	return c.JSON(helpers.CreatePaginationResponse(hydrateFeed(rows, viewerID), total, p.Offset, p.Limit))
}
