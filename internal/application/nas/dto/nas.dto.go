package dto

type CreateNASRequest struct {
	NASName         string `json:"nasname" binding:"required,max=128"`
	ShortName       string `json:"shortname" binding:"omitempty,max=32"`
	Type            string `json:"type" binding:"omitempty,max=30"`
	Ports           *int   `json:"ports" binding:"omitempty"`
	Secret          string `json:"secret" binding:"required,max=60"`
	Server          string `json:"server" binding:"omitempty,max=64"`
	Community       string `json:"community" binding:"omitempty,max=50"`
	Description     string `json:"description" binding:"omitempty,max=200"`
	RequireMa       string `json:"require_ma" binding:"omitempty,max=4"`
	LimitProxyState string `json:"limit_proxy_state" binding:"omitempty,max=4"`
}

type UpdateNASRequest struct {
	NASName         string `json:"nasname" binding:"omitempty,max=128"`
	ShortName       string `json:"shortname" binding:"omitempty,max=32"`
	Type            string `json:"type" binding:"omitempty,max=30"`
	Ports           *int   `json:"ports" binding:"omitempty"`
	Secret          string `json:"secret" binding:"omitempty,max=60"`
	Server          string `json:"server" binding:"omitempty,max=64"`
	Community       string `json:"community" binding:"omitempty,max=50"`
	Description     string `json:"description" binding:"omitempty,max=200"`
	RequireMa       string `json:"require_ma" binding:"omitempty,max=4"`
	LimitProxyState string `json:"limit_proxy_state" binding:"omitempty,max=4"`
}

type NASResponse struct {
	ID              uint   `json:"id"`
	NASName         string `json:"nasname"`
	ShortName       string `json:"shortname"`
	Type            string `json:"type"`
	Ports           *int   `json:"ports"`
	Secret          string `json:"secret"`
	Server          string `json:"server"`
	Community       string `json:"community"`
	Description     string `json:"description"`
	RequireMa       string `json:"require_ma"`
	LimitProxyState string `json:"limit_proxy_state"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type ListNASResponse struct {
	Data      []NASResponse `json:"data"`
	Total     int64         `json:"total"`
	Page      int           `json:"page"`
	PageSize  int           `json:"page_size"`
	TotalPage int           `json:"total_page"`
}

type NASFilter struct {
	NASName     string `json:"nasname" form:"nasname"`
	ShortName   string `json:"shortname" form:"shortname"`
	Type        string `json:"type" form:"type"`
	Description string `json:"description" form:"description"`
	Page        int    `json:"page" form:"page" binding:"min=1"`
	PageSize    int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}
