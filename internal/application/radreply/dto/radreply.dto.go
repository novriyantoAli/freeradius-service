package dto

type CreateRadreplyRequest struct {
	Username  string `json:"username" binding:"required"`
	Attribute string `json:"attribute" binding:"required"`
	Op        string `json:"op" binding:"required"`
	Value     string `json:"value" binding:"required"`
}

type UpdateRadreplyRequest struct {
	Username  string `json:"username"`
	Attribute string `json:"attribute"`
	Op        string `json:"op"`
	Value     string `json:"value"`
}

type RadreplyResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Attribute string `json:"attribute"`
	Op        string `json:"op"`
	Value     string `json:"value"`
}

type ListRadreplyResponse struct {
	Data      []RadreplyResponse `json:"data"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	TotalPage int                `json:"total_page"`
}

type RadreplyFilter struct {
	Username  string `form:"username"`
	Attribute string `form:"attribute"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=10"`
}
