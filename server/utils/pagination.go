package utils

import (
	"fmt"
	"net/http"
)

type PaginationRequestQuery struct {
	Page  int32 `query:"page"`
	Limit int32 `query:"limit"`
}

type PaginationResponse[T any] struct {
	Docs            []T   `json:"docs"`
	Limit           int32 `json:"limit"`
	Count           int   `json:"count"`
	CurrentPage     int32 `json:"page"`
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
}

func (pgRes *PaginationResponse[T]) ParseQuery(r *http.Request, maxLimit int32) error {
	reqQuery := PaginationRequestQuery{}
	if err := parseQueries(r, &reqQuery); err != nil {
		fmt.Println(err)
	}

	var pageNo int32
	if reqQuery.Page == 0 {
		pageNo = 1
	} else {
		pageNo = reqQuery.Page
	}

	var pageLimit int32
	if reqQuery.Limit == 0 {
		pageLimit = 10
	} else {
		pageLimit = reqQuery.Limit
	}
	pageLimit = min(pageLimit, maxLimit)

	pgRes.Limit = pageLimit
	pgRes.CurrentPage = pageNo
	return nil
}

func (pgRes *PaginationResponse[T]) BuildPaginationResponse() {
	pgRes.Count = len(pgRes.Docs)
	if len(pgRes.Docs) == int(pgRes.Limit+1) {
		pgRes.HasNextPage = true
		pgRes.Docs = pgRes.Docs[:len(pgRes.Docs)-1]
	}

	pgRes.HasPreviousPage = pgRes.CurrentPage > 1
	if pgRes.Docs == nil {
		pgRes.Docs = []T{}
	}
}
