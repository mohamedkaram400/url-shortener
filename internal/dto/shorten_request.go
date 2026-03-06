package dto


type ShortenUrlRequest struct {
	LongUrl 		string 		`json:"long_url" binding:"required,min=5"`
	ExpirationDays 	int     	`json:"expiration_days" binding:"required"`
	Status 			string 		`json:"status" binding:"oneof=Active Inactive"`
	CustomAlias     *string 	`json:"custom_alias" binding:"omitempty,min=5"`
}
	