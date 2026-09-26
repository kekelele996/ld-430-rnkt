package dto

import "time"

// CreateShareLinkRequest is the payload for creating a collection share link.
type CreateShareLinkRequest struct {
	ExpiresInHours int `json:"expires_in_hours" validate:"required,min=1,max=720"`
}

// ShareLinkCreatedResponse is returned once when a share link is created.
// The raw token appears only in this response; the server stores only its fingerprint.
type ShareLinkCreatedResponse struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	ShareURL  string    `json:"share_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ShareLinkResponse describes a share link without exposing the token.
type ShareLinkResponse struct {
	ID        string     `json:"id"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
}

// SharedAssetView is the public view of a published asset exposed via a share link.
type SharedAssetView struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	FileType     string   `json:"file_type"`
	FileFormat   string   `json:"file_format"`
	FileURL      string   `json:"file_url"`
	ThumbnailURL string   `json:"thumbnail_url"`
	FileSize     int64    `json:"file_size"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Tags         []string `json:"tags"`
	LicenseType  string   `json:"license_type"`
}

// SharedCollectionView is the public, published-only view of a shared collection.
type SharedCollectionView struct {
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	CoverImageURL string            `json:"cover_image_url"`
	AssetCount    int               `json:"asset_count"`
	Assets        []SharedAssetView `json:"assets"`
	ExpiresAt     time.Time         `json:"expires_at"`
}
