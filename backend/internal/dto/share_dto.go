package dto

import "time"

// CreateShareLinkRequest is the payload for creating a collection share link.
type CreateShareLinkRequest struct {
	ExpiresInHours int `json:"expires_in_hours" validate:"required,min=1,max=8760"`
}

// ShareLinkResponse describes a share link for the owner.
// Token is only populated once, at creation time; it is never stored.
type ShareLinkResponse struct {
	ID           string     `json:"id"`
	CollectionID string     `json:"collection_id"`
	Token        string     `json:"token,omitempty"`
	Status       string     `json:"status"` // active / expired / revoked
	ExpiresAt    time.Time  `json:"expires_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// SharedCollectionView is the public, unauthenticated view of a shared collection.
type SharedCollectionView struct {
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	CoverImageURL string            `json:"cover_image_url"`
	AssetCount    int               `json:"asset_count"`
	ExpiresAt     time.Time         `json:"expires_at"`
	Assets        []SharedAssetView `json:"assets"`
}

// SharedAssetView exposes only the published-asset fields that are safe for
// public sharing. It deliberately has no status, object key or uploader fields.
type SharedAssetView struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	FileType      string   `json:"file_type"`
	FileFormat    string   `json:"file_format"`
	FileURL       string   `json:"file_url"`
	ThumbnailURL  string   `json:"thumbnail_url"`
	FileSize      int64    `json:"file_size"`
	Width         int      `json:"width"`
	Height        int      `json:"height"`
	Tags          []string `json:"tags"`
	LicenseType   string   `json:"license_type"`
	DownloadCount int64    `json:"download_count"`
	ViewCount     int64    `json:"view_count"`
}
