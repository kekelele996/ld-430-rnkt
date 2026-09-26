package constants

// User-facing messages used in unified responses.
const (
	MsgOK                  = "ok"
	MsgUnauthorized        = "unauthorized"
	MsgForbidden           = "forbidden"
	MsgNotFound            = "resource not found"
	MsgValidationFailed    = "validation failed"
	MsgRateLimited         = "too many requests"
	MsgInternalError       = "internal server error"
	MsgEmailExists         = "email already registered"
	MsgInvalidCredentials  = "invalid email or password"
	MsgLoginSuccess        = "login success"
	MsgRegisterSuccess     = "register success"
	MsgUploadSuccess       = "upload success"
	MsgDownloadSuccess     = "download success"
	MsgReviewSuccess       = "review submitted"
	MsgCommercialForbidden = "commercial license assets require extra permission"
	MsgShareLinkInvalid    = "share link is invalid, expired or revoked"
)
