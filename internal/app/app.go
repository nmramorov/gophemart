package app

import (
	"context"
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
	if err != nil && err != http.ErrServerClosed {
		logger.InfoLog.Println("Shutting down jobmanager")
		a.manager.Shutdown()
		logger.ErrorLog.Fatalf("Server error: %e", err)
	}
}

func NewApp(config *config.Config) (*App, error) {
	logger.InfoLog.Printf("Application is running on addr %s", config.Address)
	logger.InfoLog.Printf("Accrual addr is %s", config.Accrual)
	logger.InfoLog.Printf("DB addr is %s", config.DatabaseURI)
	ctx := context.Background()
	cursor, err := db.GetCursor(config.DatabaseURI)
	if err != nil {
		return nil, err
	}
	manager := jobmanager.NewJobmanager(cursor, config.Accrual, &ctx)
	handler := api.NewHandler(cursor, manager)
	server := &http.Server{
		Addr:    config.Address,
		Handler: handler,
	}
	return &App{
		config:  config,
		manager: manager,
		Server:  server,
	}, nil
}
