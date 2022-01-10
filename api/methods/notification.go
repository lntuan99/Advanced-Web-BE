package methods

import (
	"advanced-web.hcmus/api/base"
	"advanced-web.hcmus/model"
	"github.com/gin-gonic/gin"
)

func MethodGetNotificationList(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	var notificationMappingArray = make([]model.UserNotificationMapping, 0)
	model.DBInstance.
		Joins("INNER JOIN notifications ON notifications.id = user_notification_mappings.notification_id").
		Where("user_notification_mappings.user_id = ?", user.ID).
		Order("notifications.created_at DESC").
		Offset(base.GetIntQuery(c, "page") * base.PageSizeLimit).
		Limit(base.PageSizeLimit).
		Preload("User").
		Preload("Notification").
		Find(&notificationMappingArray)

	var notificationResArray = make([]model.NotificationRes, 0)
	for _, mapping := range notificationMappingArray {
		notificationResArray = append(notificationResArray, mapping.Notification.ToRes())
	}

	return true, base.CodeSuccess, notificationResArray
}
