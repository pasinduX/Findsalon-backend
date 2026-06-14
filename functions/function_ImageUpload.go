package functions

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"findsalon-backend/integrations"

	"github.com/gofiber/fiber/v2"
)

var allowedImageExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

func SaveUploadedImage(c *fiber.Ctx, fieldName, subFolder string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	if err := validateImageFile(file); err != nil {
		return "", err
	}
	return UploadToS3(file, sanitizeS3Folder(subFolder))
}

func validateImageFile(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExtensions[ext] {
		return fmt.Errorf("invalid file type: %s", ext)
	}
	maxBytes := int64(integrations.MaxImageSizeMb) * 1024 * 1024
	if file.Size > maxBytes {
		return fmt.Errorf("image size exceeds %d MB", integrations.MaxImageSizeMb)
	}
	return nil
}

func sanitizeS3Folder(folder string) string {
	folder = strings.TrimSpace(folder)
	folder = strings.Trim(folder, "/")
	folder = filepath.ToSlash(folder)
	if folder == "" || folder == "." {
		return "general"
	}
	return folder
}
