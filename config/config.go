package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort			string
	BaseUrl			string
	DatabaseURL		string
	RedisServer		string
	JWTSecretKey	string	
	ShortCodeLenght	int	
	AccessTokenTime		int
	RefrashTokenTime	int
}

func LoadData() *Config {

	refrashTokenTimeStr := getOrDefault("REFRASH_TOKEN_TIME", "30")
	refrashTokenTimeInt, _ := strconv.Atoi(refrashTokenTimeStr)

	accessTokenTimeStr := getOrDefault("ACCESS_TOKEN_TIME", "15")
	accessTokenTimeInt, _ := strconv.Atoi(accessTokenTimeStr)

	shortCodeLengthStr := getOrDefault("SHORT_CODE_LENGTH", "7")
	shortCodeLengthInt, _ := strconv.Atoi(shortCodeLengthStr)

	return &Config{
		AppPort:           getOrDefault("APP_PORT", ":9000"),

		// Postgres connection data
		BaseUrl:  	      	  os.Getenv("BASE_URL"),
		DatabaseURL:  	      os.Getenv("DB_URL"),
		RedisServer:  	      os.Getenv("REDIS_SERVER"),
		JWTSecretKey:		  os.Getenv("JWT_SECRET_KEY"),
		ShortCodeLenght:      shortCodeLengthInt,
		AccessTokenTime:	  accessTokenTimeInt,
		RefrashTokenTime:	  refrashTokenTimeInt,
	}
}

func getOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
