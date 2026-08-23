package mail

import (
	"backend/database"
	"backend/models"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	queueSize           = 64
	defaultDailyLimit   = 80 // Resend ücretsiz tier ~100/gün
	perMessageTimeout   = 20 * time.Second
	statusSent          = "sent"
	statusFailed        = "failed"
	statusSkippedQuota  = "skipped_quota"
	statusLoggedToStdIO = "logged"
)

// Job, kuyruğa giren tek bir mail. Handler'da düz değerlerle
// doldurulur; Fiber context'i havuzdan yeniden kullandığı için
// goroutine'e *fiber.Ctx taşınmaz.
type Job struct {
	To       string
	Subject  string
	Template string
	Data     templateData
}

var (
	queue     chan Job
	startOnce sync.Once
)

// StartWorker, tek tüketici goroutine'i başlatır. Tek worker olması
// hem goroutine sayısını sınırlar hem Resend'in 2 istek/sn limitinin
// altında doğal olarak serileştirir hem de kotayı tek yerde uygular.
func StartWorker() {
	startOnce.Do(func() {
		queue = make(chan Job, queueSize)
		go func() {
			for job := range queue {
				process(job)
			}
		}()
		log.Printf("[mail] worker başlatıldı (enabled=%v, günlük limit=%d)", enabled(), dailyLimit())
	})
}

// Enqueue, maili kuyruğa atar. Kuyruk doluysa mail düşürülür —
// çağıranın isteği asla bunun yüzünden başarısız olmaz.
func Enqueue(job Job) {
	if queue == nil {
		log.Println("[mail] worker başlatılmamış, mail atlandı")
		return
	}
	if strings.TrimSpace(job.To) == "" {
		log.Printf("[mail] alıcı adresi boş, mail atlandı (%s)", job.Template)
		return
	}

	select {
	case queue <- job:
	default:
		log.Printf("[mail] kuyruk dolu, mail düşürüldü (%s -> %s)", job.Template, job.To)
	}
}

func process(job Job) {
	html, err := renderJob(job)
	if err != nil {
		log.Printf("[mail] şablon render edilemedi (%s): %v", job.Template, err)
		record(job, statusFailed, "", err.Error())
		return
	}

	// Geliştirme sürücüsü: gerçek gönderim yok, şablon stdout'a basılır.
	if !enabled() {
		log.Printf("[mail] (DEV) %s -> %s | konu: %s\n%s", job.Template, job.To, job.Subject, html)
		record(job, statusLoggedToStdIO, "", "")
		return
	}

	if overDailyLimit() {
		log.Printf("[mail] günlük limit aşıldı, mail atlandı (%s -> %s)", job.Template, job.To)
		record(job, statusSkippedQuota, "", "")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), perMessageTimeout)
	defer cancel()

	providerID, err := sendViaResend(ctx, job.To, job.Subject, html)
	if err != nil {
		log.Printf("[mail] gönderilemedi (%s -> %s): %v", job.Template, job.To, err)
		record(job, statusFailed, "", err.Error())
		return
	}

	record(job, statusSent, providerID, "")
}

func renderJob(job Job) (string, error) {
	switch job.Template {
	case TemplateAdminNewRequest:
		return render(adminNewRequestTmpl, job.Data)
	case TemplateUserApproved:
		return render(userApprovedTmpl, job.Data)
	case TemplateUserRejected:
		return render(userRejectedTmpl, job.Data)
	default:
		return render(adminNewRequestTmpl, job.Data)
	}
}

// overDailyLimit, bugün gönderilmiş mail sayısını DB'den okur.
// Bellekteki sayaçtan farklı olarak restart'a dayanıklı.
func overDailyLimit() bool {
	var count int64
	err := database.DB.Model(&models.MailLog{}).
		Where("status = ? AND created_at >= date_trunc('day', now())", statusSent).
		Count(&count).Error
	if err != nil {
		// Sayamıyorsak göndermeye devam et; kota koruması log'dan görülür.
		log.Printf("[mail] günlük kota sayılamadı: %v", err)
		return false
	}
	return count >= int64(dailyLimit())
}

func record(job Job, status, providerID, errMsg string) {
	entry := models.MailLog{
		ToEmail:    job.To,
		Template:   job.Template,
		Status:     status,
		ProviderID: providerID,
		Error:      errMsg,
		CreatedAt:  time.Now(),
	}
	if err := database.DB.Create(&entry).Error; err != nil {
		log.Printf("[mail] mail_logs yazılamadı: %v", err)
	}
}

func enabled() bool {
	if strings.EqualFold(os.Getenv("MAIL_ENABLED"), "false") {
		return false
	}
	return os.Getenv("RESEND_API_KEY") != ""
}

func dailyLimit() int {
	if raw := os.Getenv("MAIL_DAILY_LIMIT"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			return n
		}
	}
	return defaultDailyLimit
}
