package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateSalonApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)
	salonId := c.Query("SalonId")
	if salonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required")
	}

	if isOwner, _ := functions.IsSalonOwner(userId, salonId); !isOwner {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
	}

	var update bson.M
	if err := c.BodyParser(&update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	for _, field := range []string{"SalonId", "OwnerId", "CreatedAt", "Deleted"} {
		delete(update, field)
	}

	// If the client is updating Location coordinates, keep GeoLocation in sync.
	if loc, ok := update["Location"].(bson.M); ok {
		lat, hasLat := loc["Latitude"].(float64)
		lng, hasLng := loc["Longitude"].(float64)
		if hasLat && hasLng && (lat != 0 || lng != 0) {
			geo := dto.NewGeoPoint(lat, lng)
			update["GeoLocation"] = bson.M{"type": geo.Type, "coordinates": geo.Coordinates}
		}
	}

	if err := dao.UpdateSalon(salonId, update); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update salon")
	}
	return utils.SendSuccessResponse(c)
}
