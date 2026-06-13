package api

import (
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createBarberRequest struct {
	SalonId     string   `json:"SalonId"`
	UserId      string   `json:"UserId"`
	Email       string   `json:"Email"`
	Name        string   `json:"Name"`
	Specialties []string `json:"Specialties"`
	Bio         string   `json:"Bio"`
	ImageUrl    string   `json:"ImageUrl"`
}

func CreateBarberApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var req createBarberRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if req.SalonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required")
	}
	if req.Name == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Name is required")
	}

	// Only enforce ownership when an authenticated userId is available.
	if userId != "" {
		if isOwner, _ := functions.IsSalonOwner(userId, req.SalonId); !isOwner {
			return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied: you are not the owner of this salon")
		}
	}

	// Resolve UserId from Users collection when not provided directly.
	linkedUserId := req.UserId
	if linkedUserId == "" && req.Email != "" {
		if user, err := dao.FindUserByEmail(req.Email); err == nil {
			linkedUserId = user.UserId
		}
	}

	// Prevent adding the same user twice to the same salon.
	if linkedUserId != "" {
		existing, err := dao.FindBarberByUserId(linkedUserId)
		if err == nil && existing.SalonId == req.SalonId && !existing.Deleted {
			return utils.SendErrorResponse(c, fiber.StatusConflict, "This user is already a barber at this salon")
		}
	}

	now := time.Now()
	barber := dto.Barber{
		BarberId:  uuid.New().String(),
		SalonId:   req.SalonId,
		UserId:    linkedUserId,
		Name:      req.Name,
		Specialties: req.Specialties,
		Bio:       req.Bio,
		ImageUrl:  req.ImageUrl,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
		Deleted:   false,
	}

	if err := dao.CreateBarber(barber); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create barber")
	}

	// Assign barber UserRole so dashboard routing picks it up.
	if linkedUserId != "" {
		go dao.CreateUserRole(dto.UserRole{
			RoleId:    uuid.New().String(),
			UserId:    linkedUserId,
			Role:      dto.RoleBarber,
			SalonId:   req.SalonId,
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   false,
		})
	}

	return utils.SendSuccessResponse(c)
}
