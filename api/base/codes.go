package base

const (
	// General
	CodeSuccess                     = "SUCCESS"
	CodeBadRequest                  = "BAD_REQUEST"
	CodePermissionDenied            = "PERMISSION_DENIED"
	CodeDisabledBorrowing           = "DISABLED_BORROWING"
	CodeSystemManipulateFeatureFail = "SYSTEM_MANIPULATE_FEATURE_FAIL"
	CodeInvalidDate                 = "INVALID_DATE"
	CodeNoPermission                = "NO_PERMISSION"
	CodePhoneNotExisted             = "PHONE_NOT_EXISTED"
	CodeInvalidPhoneFormat          = "INVALID_PHONE_FORMAT"
	CodeInvalidEmailFormat          = "INVALID_EMAIL_FORMAT"
	CodeInvalidLanguage             = "INVALID_LANGUAGE"
	CodeInternalError               = "INTERNAL_ERROR"
	CodePhoneExisted                = "PHONE_EXISTED"
	CodeEmailExisted                = "EMAIL_EXISTED"
	CodeEmptyEmail                  = "EMPTY_EMAIL"

	// Account
	CodeUsernameExisted               = "USERNAME_EXISTED"
	CodeEmptyPassword                 = "EMPTY_PASSWORD"
	CodePasswordAndRetypeDoesNotMatch = "PASSWORD_AND_RETYPE_DOES_NOT_MATCH"
	CodeRegisterAccountFail           = "REGISTER_ACCOUNT_FAIL"
	CodeUserCodeExisted               = "USER_CODE_EXISTED"
	CodeWrongPassword                 = "WRONG_PASSWORD"
	CodeUsernameNotExisted            = "USERNAME_NOT_EXISTED"
	CodeExpiredUserAccount            = "EXPIRED_USER_ACCOUNT"
	CodeLoginAccountFail              = "LOGIN_ACCOUNT_FAIL"
	CodeNameOfUserEmpty               = "NAME_OF_USER_EMPTY"
	CodeIdentityCardExisted           = "IDENTITY_CARD_EXISTED"

	// Classroom
	CodeCreateClassroomFail  = "CREATE_CLASSROOM_FAIL"
	CodeClassroomCodeExisted = "EXISTED_CLASSROOM_CODE"
)
