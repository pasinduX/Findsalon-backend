package api

import (
	"findsalon-backend/dao"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func FindAllSalonApi(c *fiber.Ctx) error {
	filter := bson.M{}

	if ownerId := c.Query("OwnerId"); ownerId != "" {
		filter["OwnerId"] = ownerId
	}
	if cityId := c.Query("CityId"); cityId != "" {
		filter["CityId"] = cityId
	}
	if city := c.Query("City"); city != "" {
		filter["City"] = city
	}
	if districtId := c.Query("DistrictId"); districtId != "" {
		filter["DistrictId"] = districtId
	}
	if area := c.Query("Area"); area != "" {
		filter["Area"] = area
	}
	if isActive := c.Query("IsActive"); isActive != "" {
		filter["IsActive"] = isActive == "true"
	}

	salons, err := dao.FindAllSalons(filter)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve salons")
	}
	return utils.SendDataResponse(c, salons)
}
