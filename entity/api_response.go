package entity

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ApiRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"gt=0,lte=120"`
	Price       float64 `json:"price" validate:"gt=0"`
}
