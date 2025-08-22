package pagination

import (
	"fmt"

	"gorm.io/gorm"
)

type PaginatedResult[T any] struct {
	Data 		[]T 	`json:"data"`
	Page 		int 	`json:"page"`
	Limit 		int 	`json:"limit"`
	Total 		int64 	`json:"total"`
	TotalPages 	int 	`json:"total_pages"`
	Next 		string 	`json:"next"`
	Prev 		string 	`json:"prev"`
}

func Paginate[T any](db *gorm.DB, page int, limit int, baseUrl string) (PaginatedResult[T], error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var result PaginatedResult[T]
	var total int64
	var items []T

	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		return result, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	requestedPage := page
	if totalPages > 0 && requestedPage > totalPages {
		page = totalPages
	}

	offset := (page - 1) * limit

	if err := db.Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return result, err
	}

	result = PaginatedResult[T]{
		Data: 		items,
		Page: 		requestedPage,
		Limit: 		limit,
		Total: 		total,
		TotalPages: totalPages,
	}

	if requestedPage < totalPages {
		result.Next = fmt.Sprintf("%s?page=%d&limit=%d", baseUrl, requestedPage+1, limit)
	}

	if requestedPage > 1 {
		prevPage := requestedPage - 1
		if prevPage > totalPages {
			prevPage = totalPages
		}
		result.Prev = fmt.Sprintf("%s?page=%d&limit=%d", baseUrl, prevPage, limit)
	}

	return result, nil
}