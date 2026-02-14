package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type App struct {
	Name    string `yaml:"name" env-required:"true"`
	Version string `yaml:"version" env-required:"true"`
}

type Server struct {
	Host string `env:"SERVER_HOST" env-default:"localhost"`
	Port string `env:"SERVER_PORT" env-default:"80"`
}

type Config struct {
	App    App    `yaml:"app"`
	Server Server `yaml:"server"`
}

func InitConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = cleanenv.ReadConfig("config.yaml", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
