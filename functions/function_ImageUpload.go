package functions

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"findsalon-backend/integrations"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileName := uuid.New().String() + ext
	folderPath := filepath.Join(integrations.StorageBasePath, subFolder)
	if err := os.MkdirAll(folderPath, 0o755); err != nil {
		return "", err
	}
	destination := filepath.Join(folderPath, fileName)
	if err := c.SaveFile(file, destination); err != nil {
		return "", err
	}
	return filepath.ToSlash(filepath.Join(subFolder, fileName)), nil
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
