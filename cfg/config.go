package cfg

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AppEnv   string
	Host     string
	Port     string
	DbConfig DbConfig
}

type DbConfig struct {
	user     string
	password string
	name     string
	port     string
}

func LoadConfig() Config {
	local := flag.Bool("local", false, "Specify if running locally (not docker)")
	flag.Parse()

	port := os.Getenv("DB_PORT")
	if *local {
		port = os.Getenv("DB_EXTERNAL_PORT")
	}

	return Config{
		AppEnv: os.Getenv("APP_ENV"),
		Host:   os.Getenv("HOST"),
		Port:   os.Getenv("PORT"),
		DbConfig: DbConfig{
			user:     os.Getenv("DB_USER"),
			password: os.Getenv("DB_PASSWORD"),
			name:     os.Getenv("DB_NAME"),
			port:     port,
		},
	}
}

func (cfg *Config) DbUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbConfig.user,
		cfg.DbConfig.password,
		cfg.Host,
		cfg.DbConfig.port,
		cfg.DbConfig.name,
	)
}
