package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Asset is the core digital asset document.
type Asset struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	FileType      string             `bson:"file_type" json:"file_type"`
	FileFormat    string             `bson:"file_format" json:"file_format"`
	FileURL       string             `bson:"file_url" json:"file_url"`
	ThumbnailURL  string             `bson:"thumbnail_url" json:"thumbnail_url"`
	FileSize      int64              `bson:"file_size" json:"file_size"`
	Width         int                `bson:"width" json:"width"`
	Height        int                `bson:"height" json:"height"`
	Tags          []string           `bson:"tags" json:"tags"`
	CategoryID    primitive.ObjectID `bson:"category_id" json:"category_id"`
	UploaderID    primitive.ObjectID `bson:"uploader_id" json:"uploader_id"`
	LicenseType   string             `bson:"license_type" json:"license_type"`
	Status        string             `bson:"status" json:"status"`
	DownloadCount int64              `bson:"download_count" json:"download_count"`
	ViewCount     int64              `bson:"view_count" json:"view_count"`
	ObjectKey     string             `bson:"object_key" json:"object_key"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
