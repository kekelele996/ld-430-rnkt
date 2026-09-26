package client

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient wraps the minio-go client.
type MinioClient struct {
	Client *minio.Client
	Bucket string
}

// NewMinioClient connects to MinIO and ensures the bucket exists.
func NewMinioClient(ctx context.Context, endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("check minio bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("make minio bucket: %w", err)
		}
	}
	return &MinioClient{Client: client, Bucket: bucket}, nil
}
