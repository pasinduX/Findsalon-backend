package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveUploadedFile(file *multipart.FileHeader, destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(destDir, fileName)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return destPath, nil
}

func IsAllowedImageExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowed := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	return allowed[ext]
}
