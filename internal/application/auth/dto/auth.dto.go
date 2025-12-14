package dto

// AuthenticateRequest represents the RADIUS authentication request
type AuthenticateRequest struct {
	Username string `json:"username" binding:"required,max=64"`
	Password string `json:"password" binding:"required,max=253"`
}

// AuthenticateResponse represents the RADIUS authentication response
type AuthenticateResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	User    UserAuth    `json:"user,omitempty"`
	Replies []ReplyAttr `json:"replies,omitempty"`
}

// UserAuth represents authenticated user information
type UserAuth struct {
	Username   string      `json:"username"`
	Attributes []AttrValue `json:"attributes,omitempty"`
}

// AttrValue represents a single attribute value
type AttrValue struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

// ReplyAttr represents a reply attribute
type ReplyAttr struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

// AuthStatusResponse represents the authentication status
type AuthStatusResponse struct {
	Status string `json:"status"`
	User   string `json:"user,omitempty"`
	Error  string `json:"error,omitempty"`
}
