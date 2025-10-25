package cfg

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AppEnv string
	Host   string
	Port   string
	DbConfig DbConfig
}

type DbConfig struct {
	user     string
	password string
	Name     string
	Host     string
	Port     string
}


func LoadConfig() Config {
	internal := flag.Bool("internal", false, "internal=true if running inside docker container")
	flag.Parse()
	dbPort := os.Getenv("DB_EXTERNAL_PORT")
	dbHost := os.Getenv("DB_EXTERNAL_HOST")
	if *internal {
		dbPort = os.Getenv("DB_PORT")
		dbHost = os.Getenv("DB_HOST")
	}
	return Config{
		AppEnv: os.Getenv("APP_ENV"),
		Host:   os.Getenv("HOST"),
		Port:   os.Getenv("PORT"),
		DbConfig: DbConfig{
			user:     os.Getenv("DB_USER"),
			password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Host:     dbHost,
			Port:     dbPort,
		},
	}
}

func (cfg *Config) DbUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbConfig.user,
		cfg.DbConfig.password,
		cfg.DbConfig.Host,
		cfg.DbConfig.Port,
		cfg.DbConfig.Name,
	)
}
