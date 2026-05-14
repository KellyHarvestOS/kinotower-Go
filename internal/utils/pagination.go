package utils

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Page   int
	Size   int
	Limit  int
	Offset int
}

func NewPagination(r *http.Request) Pagination {
	page := queryInt(r, "page", 1)
	size := queryInt(r, "size", 10)
	if size > 100 {
		size = 100
	}
	return Pagination{Page: page, Size: size, Limit: size, Offset: (page - 1) * size}
}

func queryInt(r *http.Request, key string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}
