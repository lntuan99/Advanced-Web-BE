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
	CodePhoneInvalid                = "PHONE_INVALID"
	CodePhoneExisted                = "PHONE_EXISTED"
	CodeEmailExisted                = "EMAIL_EXISTED"
	CodeEmptyEmail                  = "EMPTY_EMAIL"

	// Account && User
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
	CodeUpdateUserProfileSuccess      = "UPDATE_USER_PROFILE_SUCCESS"
	CodeUpdateUserProfileFail         = "UPDATE_USER_PROFILE_FAIL"
	CodeFindUserFail                  = "FIND_USER_FAIL"
	CodeGoogleIDNotExisted            = "GOOGLE_ID_NOT_EXISTED"

	// Classroom
	CodeCreateClassroomFail                 = "CREATE_CLASSROOM_FAIL"
	CodeClassroomCodeExisted                = "EXISTED_CLASSROOM_CODE"
	CodeClassroomIDNotExisted               = "CLASSROOM_ID_NOT_EXISTED"
	CodeEmptyClassroomCode                  = "CLASSROOM_CODE_EMPTY"
	CodeEmptyClassroomName                  = "CLASSROOM_NAME_EMPTY"
	CodeInvalidClassroomInviteCode          = "INVALID_CLASSROOM_INVITE_CODE"
	CodeUserAlreadyInClassroom              = "USER_ALREADY_IN_CLASSROOM"
	CodeUserAlreadyOwnerOfClass             = "USER_ALREADY_OWNER_OF_CLASSROOM"
	CodeOnlyOwnerCanInviteMemberToClassroom = "ONLY_OWNER_CAN_INVITE_MEMBER_TO_CLASS"

	//Grade
	CodeEmptyGradeName          = "GRADE_NAME_EMPTY"
	CodeGradeUserInvalid        = "USER_IS_NOT_A_TEACHER_IN_CLASS"
	CodeGradeAlreadyInClassroom = "EXISTED_GRADE_IN_CLASSROOM"
	CodeCreateGradeFail         = "FAILED_TO_CREATE_GRADE"
	CodeGradeNotExisted         = "GRADE_NOT_EXISTED"
)
