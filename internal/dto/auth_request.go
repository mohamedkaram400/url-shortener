package dto


type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
	UserName  string `json:"username" binding:"required,min=3"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	// UserName string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}