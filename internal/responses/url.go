package responses


type URLResponse struct {
	ShortURL	 string					   `json:"short_url"`
	LongURL	 	 string					   `json:"long_url"`
}


type LinkAnalyticsResponse struct {
	ShortURL	 string		`json:"short_url"`
	LongURL	 	 string		`json:"long_url"`

	TotalClicks  int       `json:"total_clicks"`
	CustomAlias *string    `json:"custom_alias"`
	ExpiresAt   string 	   `json:"expires_at"`
	CreatedAt   string     `json:"created_at"`
	Status      string     `json:"status"`
}

