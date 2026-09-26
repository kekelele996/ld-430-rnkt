package util

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/assethub/assethub/internal/constants"
	"github.com/gin-gonic/gin"
)

// MaxUploadSize is the maximum accepted upload size in bytes (50 MB).
const MaxUploadSize = 50 << 20

var allowedExtensions = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".svg": true, ".webp": true, ".gif": true,
	".otf": true, ".ttf": true, ".woff": true, ".woff2": true,
	".mp4": true, ".webm": true, ".mov": true,
	".mp3": true, ".wav": true, ".flac": true,
	".psd": true, ".ai": true, ".fig": true, ".sketch": true,
	".zip": true, ".rar": true, ".pdf": true, ".fbx": true, ".obj": true, ".glb": true,
	".ppt": true, ".pptx": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
}

// ValidateUploadedFile checks file presence, extension and size.
func ValidateUploadedFile(c *gin.Context, field string) (*multipart.FileHeader, error) {
	fileHeader, err := c.FormFile(field)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constants.MsgValidationFailed, err)
	}
	if fileHeader.Size > MaxUploadSize {
		return nil, fmt.Errorf("file too large: max %d bytes", MaxUploadSize)
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] {
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}
	return fileHeader, nil
}
