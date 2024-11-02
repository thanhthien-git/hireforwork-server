package models

type PaginateDocs[T any] struct {
	Docs        []T   `json:"docs"`
	TotalDocs   int64 `json:"totalDocs"`
	CurrentPage int64 `json:"currentPage"`
	TotalPage   int64 `json:"totalPage"`
}
