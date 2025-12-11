package models

import "time"

type Friendship struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index:idx_user_friend,unique"`
	FriendID  uint      `json:"friend_id" gorm:"not null;index:idx_user_friend,unique"`
	Status    string    `json:"status" gorm:"type:varchar(20);not null;default:'pending'"` // pending, accepted, rejected
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User   Admin `json:"user" gorm:"foreignKey:UserID"`
	Friend Admin `json:"friend" gorm:"foreignKey:FriendID"`
}

// TableName specifies the table name
func (Friendship) TableName() string {
	return "friendships"
}
