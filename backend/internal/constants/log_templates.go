package constants

// Structured log message templates.
const (
	LogRequestStart     = "http request start"
	LogRequestDone      = "http request done"
	LogUploadAsset      = "asset uploaded"
	LogDownloadAsset    = "asset downloaded"
	LogReviewAsset      = "asset reviewed"
	LogLicenseChanged   = "asset license changed"
	LogAuditWritten     = "audit log written"
	LogAuthLogin        = "user login"
	LogAuthRegister     = "user registered"
	LogRateLimited      = "rate limit exceeded"
	LogMinioUpload      = "minio object uploaded"
	LogMinioError       = "minio operation failed"
	LogCacheSet         = "cache set"
	LogCacheHit         = "cache hit"
	LogShareLinkCreated = "share link created"
	LogShareLinkRevoked = "share link revoked"
)
