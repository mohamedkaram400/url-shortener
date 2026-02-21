package config

import "os"

type Config struct {
	AppPort	string
}

func LoadData() *Config {

	return &Config{
		AppPort:           getOrDefault("APP_PORT", ":9000"),
	}
}

func getOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}