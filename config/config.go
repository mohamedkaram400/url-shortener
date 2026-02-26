package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort	string

	DatabaseURL		string
	RedisServer		string
	AccessTokenTime		int
	RefrashTokenTime	int
	JWTSecretKey		string	
}

func LoadData() *Config {

	accessTokenTime, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TIME"))
	refrashTokenTime, _ := strconv.Atoi(os.Getenv("REFRASH_TOKEN_TIME"))

	return &Config{
		AppPort:           getOrDefault("APP_PORT", ":9000"),

		// Postgres connection data
		DatabaseURL:  	      os.Getenv("DB_URL"),
		RedisServer:  	      os.Getenv("REDIS_SERVER"),
		JWTSecretKey:		  os.Getenv("JWT_SECRET_KEY"),
		AccessTokenTime:	  accessTokenTime,
		RefrashTokenTime:	  refrashTokenTime,
	}
}

func getOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
