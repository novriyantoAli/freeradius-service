package dto

type CreateRadcheckRequest struct {
	Username  string `json:"username" binding:"required,max=64"`
	Attribute string `json:"attribute" binding:"required,max=64"`
	Op        string `json:"op" binding:"omitempty,max=2"`
	Value     string `json:"value" binding:"required,max=253"`
}

type UpdateRadcheckRequest struct {
	Username  string `json:"username" binding:"omitempty,max=64"`
	Attribute string `json:"attribute" binding:"omitempty,max=64"`
	Op        string `json:"op" binding:"omitempty,max=2"`
	Value     string `json:"value" binding:"omitempty,max=253"`
}

type RadcheckResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Attribute string `json:"attribute"`
	Op        string `json:"op"`
	Value     string `json:"value"`
}

type ListRadcheckResponse struct {
	Data      []RadcheckResponse `json:"data"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	TotalPage int                `json:"total_page"`
}

type RadcheckFilter struct {
	Username  string `json:"username" form:"username"`
	Attribute string `json:"attribute" form:"attribute"`
	Page      int    `json:"page" form:"page" binding:"min=1"`
	PageSize  int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}
