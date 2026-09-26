package constants

// AssetType enumerates supported digital asset file types.
type AssetType string

const (
	AssetTypeImage    AssetType = "Image"
	AssetTypeVector   AssetType = "Vector"
	AssetTypeFont     AssetType = "Font"
	AssetTypeTemplate AssetType = "Template"
	AssetTypeVideo    AssetType = "Video"
	AssetTypeAudio    AssetType = "Audio"
	AssetType3DModel  AssetType = "3DModel"
)

// LicenseType enumerates asset license types.
type LicenseType string

const (
	LicenseFree       LicenseType = "Free"
	LicenseCC0        LicenseType = "CC0"
	LicenseCCBY       LicenseType = "CC_BY"
	LicenseCCBYSA     LicenseType = "CC_BY_SA"
	LicenseCommercial LicenseType = "Commercial"
	LicenseExtended   LicenseType = "Extended"
)

// AssetStatus enumerates asset lifecycle states.
type AssetStatus string

const (
	AssetStatusDraft     AssetStatus = "Draft"
	AssetStatusPublished AssetStatus = "Published"
	AssetStatusArchived  AssetStatus = "Archived"
	AssetStatusFlagged   AssetStatus = "Flagged"
)

// ReviewResult enumerates review outcomes.
type ReviewResult string

const (
	ReviewApproved      ReviewResult = "Approved"
	ReviewRejected      ReviewResult = "Rejected"
	ReviewNeedsRevision ReviewResult = "NeedsRevision"
)

// DownloadPurpose enumerates allowed download purposes.
type DownloadPurpose string

const (
	DownloadPurposePersonal   DownloadPurpose = "Personal"
	DownloadPurposeCommercial DownloadPurpose = "Commercial"
	DownloadPurposeEditorial  DownloadPurpose = "Editorial"
)

// TagCategory enumerates tag categories.
type TagCategory string

const (
	TagCategoryStyle     TagCategory = "Style"
	TagCategoryColor     TagCategory = "Color"
	TagCategorySubject   TagCategory = "Subject"
	TagCategoryMood      TagCategory = "Mood"
	TagCategoryTechnical TagCategory = "Technical"
	TagCategoryOther     TagCategory = "Other"
)

// ShareLinkStatus describes the lifecycle state of a collection share link.
type ShareLinkStatus string

const (
	ShareLinkActive  ShareLinkStatus = "active"
	ShareLinkExpired ShareLinkStatus = "expired"
	ShareLinkRevoked ShareLinkStatus = "revoked"
)

// ValidAssetTypes is the set of valid AssetType values.
var ValidAssetTypes = map[AssetType]bool{
	AssetTypeImage: true, AssetTypeVector: true, AssetTypeFont: true,
	AssetTypeTemplate: true, AssetTypeVideo: true, AssetTypeAudio: true, AssetType3DModel: true,
}

// ValidLicenseTypes is the set of valid LicenseType values.
var ValidLicenseTypes = map[LicenseType]bool{
	LicenseFree: true, LicenseCC0: true, LicenseCCBY: true,
	LicenseCCBYSA: true, LicenseCommercial: true, LicenseExtended: true,
}

// ValidAssetStatuses is the set of valid AssetStatus values.
var ValidAssetStatuses = map[AssetStatus]bool{
	AssetStatusDraft: true, AssetStatusPublished: true,
	AssetStatusArchived: true, AssetStatusFlagged: true,
}

// ValidReviewResults is the set of valid ReviewResult values.
var ValidReviewResults = map[ReviewResult]bool{
	ReviewApproved: true, ReviewRejected: true, ReviewNeedsRevision: true,
}
