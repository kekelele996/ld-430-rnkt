package util

import (
	"fmt"
	"strings"
)

// ThumbnailSuffix is appended to the object key to build the thumbnail key.
const ThumbnailSuffix = "_thumb"

// ThumbnailKey derives a thumbnail object key from the original object key.
func ThumbnailKey(objectKey string) string {
	ext := strings.LastIndex(objectKey, ".")
	if ext == -1 {
		return objectKey + ThumbnailSuffix
	}
	return objectKey[:ext] + ThumbnailSuffix + objectKey[ext:]
}

// PublicURL builds a public URL for an object in a bucket.
func PublicURL(endpoint, bucket, objectKey string, useSSL bool) string {
	scheme := "http"
	if useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, endpoint, bucket, objectKey)
}
