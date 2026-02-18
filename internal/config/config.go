package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/numbergroup/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	Storage    string `yaml:"storage" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:9091"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func LoadConfig() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("config path if not exists")
	}
	if _, err := os.Stat(cfgPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("config file is not exists")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatal(err.Error())
	}

	return &cfg
}
