package mail

import (
	"backend/security"
	"fmt"
	"log"
	"os"
	"strings"
)

// Şablon adları; mail_logs.template alanına da bu değerler yazılır.
const (
	TemplateAdminNewRequest = "book_request_admin"
	TemplateUserApproved    = "book_request_approved"
	TemplateUserRejected    = "book_request_rejected"
	TemplateVerifyEmail     = "verify_email"
	TemplatePasswordReset   = "password_reset"
)

// BookRequestInfo, üç bildirimin de ihtiyaç duyduğu düz veri.
type BookRequestInfo struct {
	Username      string
	UserEmail     string
	Title         string
	Authors       string
	ISBN          string
	Publisher     string
	PublishDate   string
	Pages         int
	CoverURL      string
	UserNote      string
	AdminNote     string
	BookSlug      string
	MetadataFound bool
}

type AccountMailInfo struct {
	Username string
	Email    string
	URL      string
}

func NotifyVerifyEmail(info AccountMailInfo) {
	if !validRecipient(info.Email, TemplateVerifyEmail) {
		return
	}
	Enqueue(Job{
		To:       info.Email,
		Subject:  "MihrimahSiir e-posta doğrulama",
		Template: TemplateVerifyEmail,
		Data: templateData{
			Username:  info.Username,
			ActionURL: info.URL,
		},
	})
}

func NotifyPasswordReset(info AccountMailInfo) {
	if !validRecipient(info.Email, TemplatePasswordReset) {
		return
	}
	Enqueue(Job{
		To:       info.Email,
		Subject:  "MihrimahSiir şifre sıfırlama",
		Template: TemplatePasswordReset,
		Data: templateData{
			Username:  info.Username,
			ActionURL: info.URL,
		},
	})
}

// NotifyAdminNewBookRequest, yeni talep bildirimini admin'e kuyruğa atar.
func NotifyAdminNewBookRequest(info BookRequestInfo) {
	to := os.Getenv("ADMIN_NOTIFY_EMAIL")
	if !validRecipient(to, TemplateAdminNewRequest) {
		return
	}

	Enqueue(Job{
		To:       to,
		Subject:  fmt.Sprintf("Yeni kitap talebi: %s", fallbackTitle(info.Title, info.ISBN)),
		Template: TemplateAdminNewRequest,
		Data: templateData{
			Username:      info.Username,
			Title:         info.Title,
			Authors:       info.Authors,
			ISBN:          info.ISBN,
			Publisher:     info.Publisher,
			PublishDate:   info.PublishDate,
			Pages:         info.Pages,
			CoverURL:      info.CoverURL,
			UserNote:      info.UserNote,
			MetadataFound: info.MetadataFound,
			PanelURL:      publicURL("/book-request-management"),
		},
	})
}

// NotifyUserBookRequestApproved, onay bildirimini kullanıcıya kuyruğa atar.
func NotifyUserBookRequestApproved(info BookRequestInfo) {
	if !validRecipient(info.UserEmail, TemplateUserApproved) {
		return
	}

	bookURL := ""
	if info.BookSlug != "" {
		bookURL = publicURL("/book/" + info.BookSlug)
	}

	Enqueue(Job{
		To:       info.UserEmail,
		Subject:  fmt.Sprintf("Kitap talebin onaylandı: %s", fallbackTitle(info.Title, info.ISBN)),
		Template: TemplateUserApproved,
		Data: templateData{
			Username:    info.Username,
			Title:       info.Title,
			Authors:     info.Authors,
			ISBN:        info.ISBN,
			Publisher:   info.Publisher,
			PublishDate: info.PublishDate,
			Pages:       info.Pages,
			CoverURL:    info.CoverURL,
			BookURL:     bookURL,
		},
	})
}

// NotifyUserBookRequestRejected, red bildirimini kullanıcıya kuyruğa atar.
func NotifyUserBookRequestRejected(info BookRequestInfo) {
	if !validRecipient(info.UserEmail, TemplateUserRejected) {
		return
	}

	Enqueue(Job{
		To:       info.UserEmail,
		Subject:  fmt.Sprintf("Kitap talebin hakkında: %s", fallbackTitle(info.Title, info.ISBN)),
		Template: TemplateUserRejected,
		Data: templateData{
			Username:  info.Username,
			Title:     info.Title,
			Authors:   info.Authors,
			ISBN:      info.ISBN,
			AdminNote: info.AdminNote,
		},
	})
}

// validRecipient — Admin.Email unique ama not null değil, eski
// kayıtlarda boş olabilir.
func validRecipient(email, template string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		log.Printf("[mail] alıcı adresi yok, %s atlandı", template)
		return false
	}
	if err := security.NewValidator().ValidateEmail(email); err != nil {
		log.Printf("[mail] geçersiz alıcı adresi (%s), %s atlandı", email, template)
		return false
	}
	return true
}

func fallbackTitle(title, isbn string) string {
	if strings.TrimSpace(title) != "" {
		return title
	}
	return isbn
}

func publicURL(path string) string {
	base := strings.TrimRight(os.Getenv("APP_PUBLIC_URL"), "/")
	if base == "" {
		return ""
	}
	return base + path
}
