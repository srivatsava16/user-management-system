package shared

import (
	"user-management-system/internal/constants"

	"gorm.io/gorm"
)

// PaginationRequest contains pagination parameters
type PaginationRequest struct {
	Page     int `form:"page" binding:"min=0"`
	PageSize int `form:"page_size" binding:"min=0,max=100"`
}

// PaginationResponse contains pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResult wraps data with pagination info
type PaginatedResult struct {
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// NormalizePagination ensures valid pagination parameters
func NormalizePagination(req PaginationRequest) PaginationRequest {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = constants.DefaultPageSize
	}
	if req.PageSize > constants.MaxPageSize {
		req.PageSize = constants.MaxPageSize
	}
	return req
}

// Paginate applies pagination to a GORM query
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// CalculateTotalPages calculates total pages from total items and page size
func CalculateTotalPages(totalItems int64, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}
	return totalPages
}

// CreatePaginationResponse creates a pagination response
func CreatePaginationResponse(page, pageSize int, totalItems int64) PaginationResponse {
	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: CalculateTotalPages(totalItems, pageSize),
	}
}
