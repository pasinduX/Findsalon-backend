package api

import (
	"encoding/json"
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func DownloadProfileApi(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	profiles, _, err := dao.FindAllProfiles(authHeader, 0, 10000)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch profiles")
	}

	cleaned := []map[string]interface{}{}
	for _, profile := range profiles {
		record := map[string]interface{}{
			"UserId":          profile.UserId,
			"FullName":        profile.FullName,
			"Email":           profile.Email,
			"Phone":           profile.Phone,
			"AvatarUrl":       profile.AvatarUrl,
			"GoogleAvatarUrl": profile.GoogleAvatarUrl,
			"Provider":        profile.Provider,
			"Role":            profile.Role,
			"IsActive":        profile.IsActive,
		}
		cleaned = append(cleaned, record)
	}

	content, err := json.MarshalIndent(cleaned, "", "  ")
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "failed to marshal profile data")
	}

	filename := time.Now().Format("20060102")
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	c.Set(fiber.HeaderContentDisposition, "attachment; filename=users_export_"+filename+".json")
	return c.Status(fiber.StatusOK).Send(content)
}
