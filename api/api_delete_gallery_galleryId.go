package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func DeleteGalleryApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	galleryId := c.Query("GalleryId")
	salonId := c.Query("SalonId")
	barberId := c.Query("BarberId")

	if galleryId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "GalleryId is required")
	}
	if salonId == "" && barberId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId or BarberId is required")
	}

	ownerOk, _ := functions.IsSalonOwner(userId, salonId)
	isSalonOwner := salonId != "" && ownerOk
	barberOk, _ := functions.IsBarberUser(userId)
	isBarber := barberId != "" && barberOk

	if !isSalonOwner && !isBarber {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied")
	}

	if err := dao.DeleteGallery(galleryId); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete gallery item")
	}
	return utils.SendSuccessResponse(c)
}
