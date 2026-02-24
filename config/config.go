package config

import "os"

type Config struct {
	AppPort	string

	DatabaseURL		string
	RedisServer		string
}

func LoadData() *Config {

	return &Config{
		AppPort:           getOrDefault("APP_PORT", ":9000"),

		// Postgres connection data
		DatabaseURL:  	      os.Getenv("DB_URL"),
		RedisServer:  	      os.Getenv("REDIS_SERVER"),
	}
}

func getOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

