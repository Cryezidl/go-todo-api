package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Port string `env:"PORT" env-default:"8080"`
		Env  string `env:"APP_ENV" env-default:"local"`
	}

	HTTP struct {
		ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
		WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	}

	DB struct {
		Port     string `env:"DB_PORT" env-default:"5432"`
		Host     string `env:"DB_HOST" env-required:"true"`
		User     string `env:"DB_USER" env-required:"true"`
		Password string `env:"DB_PASSWORD" env-required:"true"`
		Name     string `env:"DB_NAME" env-required:"true"`
		SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
	}

	JWT struct {
		Secret     string        `env:"JWT_SECRET" env-required:"true"`
		Expiration time.Duration `env:"JWT_EXPIRATION_HOURS" env-default:"24h"`
	}
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Fatalf("config error: %v", err)
	}
	return &cfg
}
