package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port      string   `json:"port" env-required:"true"`
	Endpoints []string `json:"endpoints" env-required:"true"`
}

func MustRead(cfgpath string) *Config {
	var cfg Config

	err := cleanenv.ReadConfig(cfgpath, &cfg)
	if err != nil {
		log.Fatalf("Error while reading config file: %v", err)
	}

	return &cfg
}
