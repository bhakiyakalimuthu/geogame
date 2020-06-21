package config

import (
	"os"
	"time"
)

const (
	EnvDev     = "dev"
	EnvProd    = "prod"
)

type Config struct {
	Env string `env:"ENV" envDefault:"dev"`
	HTTPTimeOut int `env:"HTTP_TIME_OUT" envDefault:"30"`
	DBTimeOut time.Duration `env:"DB_TIME_OUT" envDefault:"3"`
	TokenSecret string `json:"TOKEN_SECRET"`
}

func NewConfig()  *Config{
	return &Config{
		Env:         EnvDev,
		HTTPTimeOut: 30,
		DBTimeOut: time.Second*3,
		TokenSecret:os.Getenv("secret"),
	}
}



