package utils

import "math"

type Pagination struct {
	Items       any   `json:"items"`
	Total       int64 `json:"total"`
	Limit       int64 `json:"limit"`
	CurrentPage int64 `json:"current_page,omitempty"`
	NextPage    int64 `json:"next_page,omitempty"`
	PrevPage    int64 `json:"prev_page,omitempty"`
	TotalPage   int64 `json:"total_page,omitempty"`
}

func GetOffset(page, limit int64) int64 {
	return (page - 1) * limit
}

func PaginateEmpty() Pagination {
	return Pagination{
		Items:       []any{},
		Total:       0,
		Limit:       0,
		CurrentPage: 1,
		NextPage:    1,
		PrevPage:    1,
		TotalPage:   1,
	}
}

func PaginatePageLimit(data any, total, page, limit int64) Pagination {

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	prevPage := 1
	if page-1 > 0 {
		prevPage = int(page) - 1
	}

	nextPage := page
	if page+1 <= int64(totalPage) {
		nextPage = page + 1
	}

	return Pagination{
		Items:       data,
		Total:       total,
		Limit:       limit,
		CurrentPage: page,
		PrevPage:    int64(prevPage),
		NextPage:    nextPage,
		TotalPage:   int64(totalPage),
	}
}

func PaginateOffsetLimit(data any, total, offset, limit int64) Pagination {
	var page int
	if offset >= total {
		page = -1
	} else {
		page = int(offset/limit) + 1
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	prevPage := 1
	if page-1 > 0 {
		prevPage = page - 1
	}

	nextPage := page
	if page+1 <= totalPage {
		nextPage = page + 1
	}

	return Pagination{
		Items:       data,
		Total:       total,
		Limit:       limit,
		CurrentPage: int64(page),
		PrevPage:    int64(prevPage),
		NextPage:    int64(nextPage),
		TotalPage:   int64(totalPage),
	}
}
