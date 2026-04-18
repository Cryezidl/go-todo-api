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

	_ = cleanenv.ReadConfig(".env", &cfg)

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Printf("Warning: failed to read env: %v", err)
	}

	if cfg.DB.Host == "" {
		log.Fatal("DB_HOST is required")
	}
	if cfg.DB.User == "" {
		log.Fatal("DB_USER is required")
	}
	if cfg.DB.Password == "" {
		log.Fatal("DB_PASSWORD is required")
	}
	if cfg.DB.Name == "" {
		log.Fatal("DB_NAME is required")
	}
	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	log.Printf("Config loaded: DB_HOST=%s, DB_NAME=%s, APP_ENV=%s",
		cfg.DB.Host, cfg.DB.Name, cfg.App.Env)

	return &cfg
}
