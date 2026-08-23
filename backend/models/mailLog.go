package models

import "time"

// MailLog, gönderilen/atlanan mailleri kaydeder. İki işi var:
// günlük kotayı restart'a dayanıklı şekilde saymak ve denetim izi bırakmak.
type MailLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ToEmail    string    `json:"to_email" gorm:"type:varchar(255)"`
	Template   string    `json:"template" gorm:"type:varchar(64)"`
	Status     string    `json:"status" gorm:"type:varchar(24);index"` // sent, failed, skipped_quota, logged
	ProviderID string    `json:"provider_id" gorm:"type:varchar(64)"`
	Error      string    `json:"error" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at" gorm:"index"`
}

func (MailLog) TableName() string {
	return "mail_logs"
}
