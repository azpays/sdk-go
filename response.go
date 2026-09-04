package azpays

// Response is the standard AzPays API response envelope.
// All API endpoints return this format: {"ok": bool, "msg": string, "data": T}
type Response[T any] struct {
	OK   bool   `json:"ok"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// PaginatedResponse is the AzPays paginated response envelope.
// Paginated endpoints return: {"success": bool, "message": string, "data": [...], "pagination": {...}}
type PaginatedResponse[T any] struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Data       []T            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta contains pagination metadata.
type PaginationMeta struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	From        int   `json:"from"`
	To          int   `json:"to"`
}

// ListParams is the base filter for all paginated list requests.
type ListParams struct {
	// Page number (1-indexed). Default: 1.
	Page int `json:"page,omitempty"`

	// PerPage items per page. Default: 15.
	PerPage int `json:"per_page,omitempty"`

	// Search is a free-text search query.
	Search string `json:"search,omitempty"`
}
