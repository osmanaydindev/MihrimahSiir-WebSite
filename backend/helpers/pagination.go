package helpers

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type PaginationParams struct {
	Offset int
	Limit  int
}

type PaginationResponse struct {
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Offset  int         `json:"offset"`
	Limit   int         `json:"limit"`
	HasMore bool        `json:"has_more"`
}

func GetPaginationParams(c *fiber.Ctx) PaginationParams {
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return PaginationParams{
		Offset: offset,
		Limit:  limit,
	}
}

func CreatePaginationResponse(data interface{}, total int64, offset int, limit int) PaginationResponse {
	hasMore := int64(offset+limit) < total

	return PaginationResponse{
		Data:    data,
		Total:   total,
		Offset:  offset,
		Limit:   limit,
		HasMore: hasMore,
	}
}
