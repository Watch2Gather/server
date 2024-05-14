package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"

	configs "github.com/Watch2Gather/server/pkg/config"
)

type (
	Config struct {
		configs.App  `yaml:"app"`
		PG           `yaml:"postgres"`
		configs.Log  `yaml:"logger"`
		configs.HTTP `yaml:"http"`
		configs.WS   `yaml:"ws"`
	}
	PG struct {
		DsnURL  string `env-required:"true" yaml:"dsn_url" env:"PG_DSN_URL"`
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"POOL_MAX"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = cleanenv.ReadConfig(dir+"/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %e", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
