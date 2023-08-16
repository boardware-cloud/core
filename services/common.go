package services

type Pagination struct {
	Index int64
	Limit int64
	Total int64
}

type List[T any] struct {
	Data       []T
	Pagination Pagination
}
