package api

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func FindSalonServiceApi(c *fiber.Ctx) error {
	serviceId := c.Query("ServiceId")
	if serviceId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "ServiceId is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var service dto.SalonService
	collection := dbConfig.DATABASE.Collection(dbConfig.SALONSERVICES_COLLECTION)
	err := collection.FindOne(ctx, bson.M{"ServiceId": serviceId, "Deleted": false}).Decode(&service)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Service not found")
	}
	return utils.SendDataResponse(c, service)
}
