package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DownloadRecord logs every asset download.
type DownloadRecord struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssetID        primitive.ObjectID `bson:"asset_id" json:"asset_id"`
	DownloaderID   primitive.ObjectID `bson:"downloader_id" json:"downloader_id"`
	DownloadedAt   time.Time          `bson:"downloaded_at" json:"downloaded_at"`
	Purpose        string             `bson:"purpose" json:"purpose"`
	LicenseVersion string             `bson:"license_version" json:"license_version"`
	IP             string             `bson:"ip" json:"ip"`
}
