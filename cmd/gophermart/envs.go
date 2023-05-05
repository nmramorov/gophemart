package main

import (
	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	Address     string `env:"RUN_ADDRESS,required" envDefault:"127.0.0.1:8080"`
	DatabaseURI string `env:"DATABASE_URI,required" envDefault:"127.0.0.1:5432"`
	Accrual     string `env:"ACCRUAL_SYSTEM_ADDRESS,required" envDefault:"127.0.0.1:8081"`
}

func NewEnvConfig() (*EnvConfig, error) {
	var config EnvConfig
	err := env.Parse(&config)
	if err != nil {
		ErrorLog.Fatalf("Error with env config: %e", err)
	}
	return &config, nil
}
