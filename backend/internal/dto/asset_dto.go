package dto

import "mime/multipart"

// UploadAssetRequest carries metadata for an asset upload.
type UploadAssetRequest struct {
	Title       string                `form:"title" validate:"required,min=1,max=200"`
	Description string                `form:"description"`
	FileType    string                `form:"file_type" validate:"required"`
	FileFormat  string                `form:"file_format" validate:"required"`
	Tags        []string              `form:"tags"`
	CategoryID  string                `form:"category_id" validate:"required"`
	LicenseType string                `form:"license_type" validate:"required"`
	Width       int                   `form:"width"`
	Height      int                   `form:"height"`
	File        *multipart.FileHeader `form:"file" validate:"required"`
}

// UpdateAssetRequest carries editable asset metadata.
type UpdateAssetRequest struct {
	Title       string   `json:"title" validate:"omitempty,min=1,max=200"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	CategoryID  string   `json:"category_id"`
	LicenseType string   `json:"license_type"`
	Status      string   `json:"status"`
}

// AssetQuery are search/filter parameters.
type AssetQuery struct {
	Keyword     string `form:"keyword"`
	FileType    string `form:"file_type"`
	LicenseType string `form:"license_type"`
	Status      string `form:"status"`
	Tag         string `form:"tag"`
	CategoryID  string `form:"category_id"`
	Sort        string `form:"sort"` // downloads | views | newest | oldest
	Page        string `form:"page"`
	PageSize    string `form:"page_size"`
}
