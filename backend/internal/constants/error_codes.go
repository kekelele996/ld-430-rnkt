package constants

// Unified response code values. 0 means success; non-zero are business error codes.
const (
	CodeOK               = 0
	CodeBadRequest       = 40000
	CodeUnauthorized     = 40100
	CodeForbidden        = 40300
	CodeNotFound         = 40400
	CodeConflict         = 40900
	CodeValidationFailed = 42200
	CodeRateLimited      = 42900
	CodeInternal         = 50000
	CodeMinioError       = 50701
	CodeUploadTooLarge   = 50702
	CodeUnsupportedFile  = 50703
)
