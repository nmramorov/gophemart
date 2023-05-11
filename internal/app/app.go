package app

import (
	"net/http"

	"github.com/nmramorov/gophemart/internal/api"
	config "github.com/nmramorov/gophemart/internal/configuration"
	"github.com/nmramorov/gophemart/internal/db"
	"github.com/nmramorov/gophemart/internal/jobmanager"
	"github.com/nmramorov/gophemart/internal/logger"
)

type App struct {
	config  *config.Config
	manager *jobmanager.Jobmanager
	Server  *http.Server
}

func (a *App) Run() {
	go a.manager.ManageJobs(a.config.Accrual)
	err := a.Server.ListenAndServe()
	if err != nil {
		logger.ErrorLog.Fatalf("Server error: %e", err)
	}
}

func NewApp(config *config.Config) *App {
	logger.InfoLog.Printf("Application is running on addr %s", config.Address)
	logger.InfoLog.Printf("Accrual addr is %s", config.Accrual)
	logger.InfoLog.Printf("DB addr is %s", config.DatabaseURI)
	handler := api.NewHandler(config.Accrual, db.GetCursor(config.DatabaseURI))
	server := &http.Server{
		Addr:    config.Address,
		Handler: handler,
	}
	return &App{
		config:  config,
		manager: handler.Manager,
		Server:  server,
	}
}
