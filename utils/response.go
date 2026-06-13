package utils

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func SendSuccessResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"code":   fiber.StatusOK,
	})
}

func SendErrorResponse(c *fiber.Ctx, statusCode int, errorMessage string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "error",
		"code":    statusCode,
		"message": errorMessage,
	})
}

func SendDataResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"code":   fiber.StatusOK,
		"data":   data,
	})
}

func SendJSONFileDownload(c *fiber.Ctx, data interface{}, filename string) error {
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to serialize data")
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	c.Set(fiber.HeaderContentDisposition, "attachment; filename="+filename+".json")
	return c.Status(fiber.StatusOK).Send(raw)
}
