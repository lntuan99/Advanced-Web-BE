package model

import "time"

type Notification struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	Title     string
	Message   string
	Payload   string
	Users     []User `gorm:"many2many:user_notification_mappings"`
}

type NotificationRes struct {
	ID        uint   `json:"id"`
	CreatedAt int64  `json:"createdAt"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Payload   string `json:"payload"`
}

func (notification Notification) ToRes() NotificationRes {
	return NotificationRes{
		ID:        notification.ID,
		CreatedAt: notification.CreatedAt.Unix(),
		Title:     notification.Title,
		Message:   notification.Message,
		Payload:   notification.Payload,
	}
}
