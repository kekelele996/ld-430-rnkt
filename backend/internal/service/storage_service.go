package service

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/assethub/assethub/internal/client"
	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/util"
	"github.com/minio/minio-go/v7"
)

// StorageService wraps MinIO object operations.
type StorageService struct {
	minio  *client.MinioClient
	logger *slog.Logger
}

// NewStorageService creates a StorageService.
func NewStorageService(minio *client.MinioClient, logger *slog.Logger) *StorageService {
	return &StorageService{minio: minio, logger: logger}
}

// UploadObject stores a file and returns the object key, public URL and thumbnail URL.
func (s *StorageService) UploadObject(ctx context.Context, fileHeader *multipart.FileHeader, objectKey string) (string, string, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", "", "", fmt.Errorf("open upload file: %w", err)
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	info, err := s.minio.Client.PutObject(ctx, s.minio.Bucket, objectKey, file, fileHeader.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", "", "", fmt.Errorf("put object: %w", err)
	}
	s.logger.Info(constants.LogMinioUpload, "bucket", s.minio.Bucket, "object", info.Key, "size", info.Size)

	publicURL := util.PublicURL(s.minio.Client.EndpointURL().Host, s.minio.Bucket, objectKey, s.minio.Client.EndpointURL().Scheme == "https")
	thumbnailKey := util.ThumbnailKey(objectKey)
	thumbnailURL := util.PublicURL(s.minio.Client.EndpointURL().Host, s.minio.Bucket, thumbnailKey, s.minio.Client.EndpointURL().Scheme == "https")
	return objectKey, publicURL, thumbnailURL, nil
}

// GenObjectKey builds a unique object key for an upload.
func (s *StorageService) GenObjectKey(userID, filename string) string {
	ext := ""
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			ext = filename[i:]
			break
		}
	}
	return fmt.Sprintf("uploads/%s/%d%s", userID, time.Now().UnixNano(), ext)
}
