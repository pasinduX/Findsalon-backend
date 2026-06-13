package functions

import (
	"math"
	"strconv"

	"findsalon-backend/dto"
	"findsalon-backend/integrations"

	"github.com/gofiber/fiber/v2"
)

type PaginationParams struct {
	Page     int
	PageSize int
}

func GetPaginationParams(c *fiber.Ctx) PaginationParams {
	page := 1
	pageSize := integrations.DefaultPageSize

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("pageSize"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize < 1 {
		pageSize = integrations.DefaultPageSize
	}
	return PaginationParams{Page: page, PageSize: pageSize}
}

func BuildSkipLimit(p PaginationParams) (int64, int64) {
	return int64((p.Page - 1) * p.PageSize), int64(p.PageSize)
}

func BuildPaginatedResponse(data interface{}, total int, p PaginationParams) dto.PaginatedResponse {
	totalPages := int(math.Ceil(float64(total) / float64(p.PageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	return dto.PaginatedResponse{
		Data:       data,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalCount: total,
		TotalPages: totalPages,
	}
}
