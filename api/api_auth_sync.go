package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type authSyncRequest struct {
	Email     string `json:"Email"`
	FullName  string `json:"FullName"`
	Provider  string `json:"Provider"`
	GoogleId  string `json:"GoogleId"`
	AvatarUrl string `json:"AvatarUrl"`
}

// AuthSyncApi finds the user by email and creates them if they don't exist yet.
// Called by the frontend immediately after Auth0 login to ensure a DB record exists.
func AuthSyncApi(c *fiber.Ctx) error {
	var req authSyncRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if req.Email == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Email is required")
	}

	existing, err := dao.FindUserByEmail(req.Email)
	if err == nil {
		if existing.GoogleId == "" && req.GoogleId != "" {
			update := bson.M{"GoogleId": req.GoogleId}
			if req.Provider != "" {
				update["Provider"] = req.Provider
			}
			if updateErr := dao.UpdateUser(existing.UserId, update); updateErr == nil {
				existing.GoogleId = req.GoogleId
				if req.Provider != "" {
					existing.Provider = req.Provider
				}
				existing.UpdatedAt = time.Now()
			}
		}
		// User already exists — return their record
		return utils.SendDataResponse(c, dto.UserResponse{
			UserId:          existing.UserId,
			FullName:        existing.FullName,
			Email:           existing.Email,
			AvatarUrl:       existing.AvatarUrl,
			GoogleAvatarUrl: existing.GoogleAvatarUrl,
			Provider:        existing.Provider,
			GoogleId:        existing.GoogleId,
			Role:            existing.Role,
			IsActive:        existing.IsActive,
			CreatedAt:       existing.CreatedAt,
			UpdatedAt:       existing.UpdatedAt,
		})
	}

	// Create new user
	now := time.Now()
	provider := req.Provider
	if provider == "" {
		provider = "auth0"
	}
	newUser := dto.User{
		UserId:          uuid.New().String(),
		FullName:        req.FullName,
		Email:           req.Email,
		Provider:        provider,
		GoogleId:        req.GoogleId,
		GoogleAvatarUrl: req.AvatarUrl,
		Role:            dto.RoleUser,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
		Deleted:         false,
	}
	if err := dao.CreateUser(newUser); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	return utils.SendDataResponse(c, dto.UserResponse{
		UserId:          newUser.UserId,
		FullName:        newUser.FullName,
		Email:           newUser.Email,
		GoogleAvatarUrl: newUser.GoogleAvatarUrl,
		Provider:        newUser.Provider,
		GoogleId:        newUser.GoogleId,
		Role:            newUser.Role,
		IsActive:        newUser.IsActive,
		CreatedAt:       newUser.CreatedAt,
		UpdatedAt:       newUser.UpdatedAt,
	})
}
