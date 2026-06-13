package api

import (
	"findsalon-backend/functions"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// UploadImageApi accepts a multipart form with:
//   - image  (required) — the image file
//   - Folder (optional) — S3 subfolder, e.g. "salon-covers", "barber-profiles", "gallery"
//
// Response: { status, code, data: { Url: "https://..." } }
func UploadImageApi(c *fiber.Ctx) error {
	folder := c.FormValue("Folder")
	if folder == "" {
		folder = "general"
	}

	url, err := functions.SaveUploadedImage(c, "image", folder)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Image upload failed: "+err.Error())
	}

	return utils.SendDataResponse(c, fiber.Map{"Url": url})
}
