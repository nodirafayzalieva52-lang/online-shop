package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPPort string `env:"HTTP_PORT" envDefault:":8080"`
	Postgres Postgres
	JWT      JWT
}

type Postgres struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName   string `env:"DB_NAME" envDefault:"shop_db"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

type JWT struct {
	Secret string `env:"JWT_SECRET" envDefault:"super-secret-key"`
	TTL    string `env:"JWT_TTL" envDefault:"24h"`
}

func New(path string) (*Config, error) {
	var conf Config

	err := cleanenv.ReadConfig(path, &conf)
	if err != nil {
		return nil, fmt.Errorf("cleanenv.ReadConfig: %w", err)
	}

	return &conf, nil
}
