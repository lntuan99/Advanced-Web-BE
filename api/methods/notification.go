package methods

import (
	"advanced-web.hcmus/api/base"
	req_res "advanced-web.hcmus/api/req_res_struct"
	"advanced-web.hcmus/model"
	"github.com/gin-gonic/gin"
	"strings"
)

func MethodGetNotificationList(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	dbInstance := model.DBInstance

	isRead := strings.ToLower(c.Query("is-read"))
	if isRead == "true" || isRead == "false" {
		dbInstance = dbInstance.Where("user_notification_mappings.is_read = ?", isRead)
	}

	var notificationMappingArray = make([]model.UserNotificationMapping, 0)
	dbInstance.
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

func MethodMarkReadNotification(c *gin.Context) (bool, string, interface{}) {
	userObj, _ := c.Get("user")
	user := userObj.(model.User)

	var notificationIDInfo req_res.PostMarkReadNotification
	if err := c.ShouldBind(&notificationIDInfo); err != nil {
		return false, base.CodeBadRequest, nil
	}

	var dbNotification model.Notification
	model.DBInstance.First(&dbNotification, notificationIDInfo.NotificationID)
	if dbNotification.ID == 0 {
		return false, base.CodeNotificationIDNotExisted, nil
	}

	var dbNotificationMapping model.UserNotificationMapping
	model.DBInstance.First(&dbNotificationMapping, "user_id = ? AND notification_id = ?", user.ID, notificationIDInfo.NotificationID)

	if err := model.DBInstance.
		Model(&dbNotificationMapping).
		Updates(map[string]interface{}{
			"is_read": true,
		}).Error; err != nil {
		return false, base.CodeMarkReadNotificationFail, nil
	}

	if dbNotificationMapping.ID == 0 {
		return false, base.CodeUserNotReceiveThisNotificationID, nil
	}

	return true, base.CodeSuccess, nil
}
