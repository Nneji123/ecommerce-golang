package middleware

type ErrorResponse struct {
	Error string `json:"error"`
}

type PaginationQuery struct {
	Page    int    `query:"page"`
	Limit   int    `query:"limit"`
	SortBy  string `query:"sort_by"`
	SortDir string `query:"sort_dir"`
}

type PaginatedResponse struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	Limit       int   `json:"limit"`
}
