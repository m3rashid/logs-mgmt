package utils

import (
	"fmt"
	"net/http"
)

type PaginationQuery struct {
	Page  int32 `query:"page"`
	Limit int32 `query:"limit"`
}

type PaginationResponse[T any] struct {
	Docs            []T   `json:"docs"`
	Limit           int32 `json:"limit"`
	Count           int32 `json:"count"`
	CurrentPage     int32 `json:"page"`
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
}

// Extract pagination from the query struct
// This uses type assertion to get the embedded PaginationQuery
type paginationGetter interface {
	GetPagination() PaginationQuery
}

// ParseQueryWithParams parses pagination and additional query parameters
// queryStruct should be a pointer to your custom query struct that embeds PaginationQuery
func (pgRes *PaginationResponse[T]) ParseQueryWithParams(r *http.Request, maxLimit int32, queryStruct any) error {
	if err := parseQueries(r, queryStruct); err != nil {
		fmt.Println(err)
	}

	if pg, ok := queryStruct.(paginationGetter); ok {
		pagination := pg.GetPagination()
		return pgRes.setPaginationParams(pagination.Page, pagination.Limit, maxLimit)
	}

	// If GetPagination() method doesn't exist, return an error
	return fmt.Errorf("queryStruct must implement GetPagination() method")
}

func (pgRes *PaginationResponse[T]) ParseQuery(r *http.Request, maxLimit int32) error {
	query := PaginationQuery{}
	return pgRes.ParseQueryWithParams(r, maxLimit, &query)
}

// setPaginationParams sets the pagination parameters on the response
func (pgRes *PaginationResponse[T]) setPaginationParams(page, limit, maxLimit int32) error {
	var pageNo int32
	if page == 0 {
		pageNo = 1
	} else {
		pageNo = page
	}

	var pageLimit int32
	if limit == 0 {
		pageLimit = 10
	} else {
		pageLimit = limit
	}
	pageLimit = min(pageLimit, maxLimit)

	pgRes.Limit = pageLimit
	pgRes.CurrentPage = pageNo
	return nil
}

func (pgRes *PaginationResponse[T]) BuildPaginationResponse() {
	pgRes.Count = int32(len(pgRes.Docs))
	if len(pgRes.Docs) == int(pgRes.Limit+1) {
		pgRes.HasNextPage = true
		pgRes.Docs = pgRes.Docs[:len(pgRes.Docs)-1]
	}

	pgRes.HasPreviousPage = pgRes.CurrentPage > 1
	if pgRes.Docs == nil {
		pgRes.Docs = []T{}
	}
}
