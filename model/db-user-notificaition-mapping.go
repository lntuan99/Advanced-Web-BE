package model

import "github.com/jinzhu/gorm"

type UserNotificationMapping struct {
	gorm.Model
	UserID         uint `gorm:"unique_index:user_and_notification_in_mapping"`
	User           User
	NotificationID uint `gorm:"unique_index:user_and_notification_in_mapping"`
	Notification   Notification
}
