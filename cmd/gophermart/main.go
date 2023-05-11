package main

import (
	"github.com/nmramorov/gophemart/internal/app"
	config "github.com/nmramorov/gophemart/internal/configuration"
	"github.com/nmramorov/gophemart/internal/logger"
)

func main() {
	flags := config.NewCliOptions()
	envs, err := config.NewEnvConfig()
	if err != nil {
		logger.ErrorLog.Fatal(err)
	}
	app, _err := app.NewApp(config.NewConfig(flags, envs))
	if _err != nil {
		logger.ErrorLog.Fatal(_err)
	}
	app.Run()
}
