package dto


type ShortenUrlRequest struct {
	LongUrl 		string 		`json:"long_url" binding:"required,min=5"`
	ExpirationDays 	uint     	`json:"expiration_days" binding:"required"`
	UserId  		uint    	`json:"user_id" binding:"required"`
	Status  		string    	`json:"status" binding:"required"`
	CustomAlias		string		`json:"custom_alias" binding:"required"`
}
