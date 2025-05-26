package utility

import (
    "github.com/gofiber/fiber/v2"
    "math"
)

type PaginationQuery struct {
    Page   int    `json:"page" query:"page"`
    Limit  int    `json:"limit" query:"limit"`
    Search string `json:"search" query:"search"`
}

type PaginationResponse struct {
    Data       interface{} `json:"data"`
    Page       int        `json:"page"`
    Limit      int        `json:"limit"`
    TotalRows  int64      `json:"total_rows"`
    TotalPages int        `json:"total_pages"`
}

func GetPaginationParams(c *fiber.Ctx) PaginationQuery {
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    search := c.Query("search", "")

    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }

    return PaginationQuery{
        Page:   page,
        Limit:  limit,
        Search: search,
    }
}

func CreatePaginationResponse(data interface{}, totalRows int64, params PaginationQuery) PaginationResponse {
    return PaginationResponse{
        Data:       data,
        Page:       params.Page,
        Limit:      params.Limit,
        TotalRows:  totalRows,
        TotalPages: int(math.Ceil(float64(totalRows) / float64(params.Limit))),
    }
}