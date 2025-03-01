package cfg

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AppEnv string
	Host   string
	Port   string
}

func LoadConfig() Config {
	return Config{
		AppEnv: os.Getenv("APP_ENV"),
		Host:   os.Getenv("HOST"),
		Port:   os.Getenv("PORT"),
	}
}
