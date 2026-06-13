package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
)

func FindAllGalleryApi(c *fiber.Ctx) error {
	salonId := c.Query("SalonId")
	barberId := c.Query("BarberId")

	if salonId == "" && barberId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId or BarberId is required")
	}

	var gallery []dto.Gallery
	var err error

	if barberId != "" {
		gallery, err = dao.FindAllGalleryByBarberId(barberId)
	} else {
		gallery, err = dao.FindAllGalleryBySalonId(salonId)
	}

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve gallery")
	}
	return utils.SendDataResponse(c, gallery)
}
