package http

type BaseResponse struct {
	Data interface{} `json:"data,omitempty"`
}

type ResponsePagination struct {
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	Total       int64 `json:"total"`
}

type PaginatedBaseResponse struct {
	Pagination ResponsePagination `json:"pagination"`
	Data       interface{}        `json:"data,omitempty"`
}
