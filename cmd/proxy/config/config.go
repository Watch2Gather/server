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
		configs.Log  `yaml:"log"`
		GRPC         `yaml:"grpc"`
		configs.HTTP `yaml:"http"`
	}

	GRPC struct {
		UserHost string `env-required:"true" yaml:"user_host" env:"GRPC_USER_HOST"`
		UserPort int    `env-required:"true" yaml:"user_port" env:"GRPC_USER_PORT"`
		RoomHost string `env-required:"true" yaml:"room_host" env:"GRPC_ROOM_HOST"`
		RoomPort int    `env-required:"true" yaml:"room_port" env:"GRPC_ROOM_PORT"`
		MovieHost string `env-required:"true" yaml:"movie_host" env:"GRPC_MOVIE_HOST"`
		MoviePort int    `env-required:"true" yaml:"movie_port" env:"GRPC_MOVIE_PORT"`
	}
)

func NewCofig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// debug
	fmt.Println(dir)

	err = cleanenv.ReadConfig(dir+"/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
