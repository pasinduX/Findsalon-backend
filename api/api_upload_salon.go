package api

import (
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func UploadSalonApi(c *fiber.Ctx) error {
	result := fiber.Map{}

	if _, err := c.FormFile("cover_image"); err == nil {
		coverUrl, err := functions.SaveUploadedImage(c, "cover_image", "salons/covers")
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Cover image upload failed: "+err.Error())
		}
		result["CoverImageUrl"] = coverUrl
	}

	if _, err := c.FormFile("logo"); err == nil {
		logoUrl, err := functions.SaveUploadedImage(c, "logo", "salons/logos")
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Logo upload failed: "+err.Error())
		}
		result["LogoUrl"] = logoUrl
	}

	if len(result) == 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "No valid image fields provided (cover_image or logo)")
	}

	return utils.SendDataResponse(c, result)
}
