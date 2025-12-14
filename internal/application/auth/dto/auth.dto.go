package dto

// CreateAuthRequest represents a request to create authentication credentials
type CreateAuthRequest struct {
	Username   string                `json:"username" binding:"required,max=64"`
	Password   string                `json:"password" binding:"required,max=253"`
	Attributes []CreateAuthAttribute `json:"attributes" binding:"omitempty"`
	ReplyAttrs []CreateAuthAttribute `json:"reply_attributes" binding:"omitempty"`
}

// CreateAuthAttribute represents an attribute to be created
type CreateAuthAttribute struct {
	Attribute string `json:"attribute" binding:"required,max=64"`
	Value     string `json:"value" binding:"required,max=253"`
	Op        string `json:"op" binding:"omitempty,max=2"` // Default: ":=" for radcheck, "+=" for radreply
}

// CreateAuthResponse represents the response after creating authentication credentials
type CreateAuthResponse struct {
	Username   string                   `json:"username"`
	Password   string                   `json:"password"`
	Attributes []AuthCreateAttrResponse `json:"attributes"`
	ReplyAttrs []AuthCreateAttrResponse `json:"reply_attributes"`
}

// AuthCreateAttrResponse represents created attribute information
type AuthCreateAttrResponse struct {
	ID        uint   `json:"id"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
	Op        string `json:"op"`
}
