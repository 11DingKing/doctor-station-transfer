package pagination

import "strconv"

type Request struct {
	Limit, Offset int
	Query, Sort   string
}

func Parse(limit, offset string) Request {
	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)
	if l <= 0 || l > 100 {
		l = 20
	}
	if o < 0 {
		o = 0
	}
	return Request{Limit: l, Offset: o, Query: "", Sort: "created_at desc"}
}

type Result[T any] struct {
	Items  []T `json:"items"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
