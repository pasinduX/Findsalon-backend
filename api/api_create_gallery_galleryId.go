package api

import (
	"strconv"
	"strings"
	"time"

	"findsalon-backend/dao"
	"findsalon-backend/dto"
	"findsalon-backend/functions"
	"findsalon-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createGalleryRequest struct {
	SalonId   string `json:"SalonId"`
	BarberId  string `json:"BarberId"`
	ImageUrl  string `json:"ImageUrl"`
	Caption   string `json:"Caption"`
	SortOrder int    `json:"SortOrder"`
}

func CreateGalleryApi(c *fiber.Ctx) error {
	userId, _ := c.Locals("userId").(string)

	var req createGalleryRequest
	contentType := strings.ToLower(c.Get("Content-Type"))
	if strings.Contains(contentType, "application/json") {
		if err := c.BodyParser(&req); err != nil {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
		}
	} else {
		req.SalonId = c.FormValue("SalonId")
		req.BarberId = c.FormValue("BarberId")
		req.ImageUrl = c.FormValue("ImageUrl")
		req.Caption = c.FormValue("Caption")
		if soStr := c.FormValue("SortOrder"); soStr != "" {
			req.SortOrder, _ = strconv.Atoi(soStr)
		}
	}

	if req.SalonId == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "SalonId is required")
	}

	isSalonOwner, _ := functions.IsSalonOwner(userId, req.SalonId)
	isBarberUser, _ := functions.IsBarberUser(userId)
	isBarber := req.BarberId != "" && isBarberUser

	if !isSalonOwner && !isBarber {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Access denied")
	}

	imageUrl := req.ImageUrl
	if imageUrl == "" {
		var err error
		imageUrl, err = functions.SaveUploadedImage(c, "image", "gallery/"+req.SalonId)
		if err != nil {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Image upload failed: "+err.Error())
		}
	}

	// BarberId is nullable: only set it when a non-empty value was supplied.
	var barberId *string
	if req.BarberId != "" {
		b := req.BarberId
		barberId = &b
	}

	gallery := dto.Gallery{
		GalleryId: uuid.New().String(),
		SalonId:   req.SalonId,
		BarberId:  barberId,
		ImageUrl:  imageUrl,
		Caption:   req.Caption,
		SortOrder: req.SortOrder,
		CreatedAt: time.Now(),
		Deleted:   false,
	}

	if err := dao.CreateGallery(gallery); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create gallery item")
	}
	return utils.SendSuccessResponse(c)
}
