package dto

// DownloadRequest records download purpose and license version.
type DownloadRequest struct {
	Purpose        string `json:"purpose" validate:"required,oneof=Personal Commercial Editorial"`
	LicenseVersion string `json:"license_version"`
}
