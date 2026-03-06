package responses

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    interface{}    `json:"data,omitempty"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

