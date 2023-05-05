package main

import "net/http"

type App struct {
	config  *Config
	manager *Jobmanager
	Server  *http.Server
}

func (a *App) Run() {
	go a.manager.ManageJobs(a.config.Accrual)
	a.Server.ListenAndServe()
}

func NewApp(config *Config) *App {
	InfoLog.Printf("Application is running on addr %s", config.Address)
	InfoLog.Printf("Accrual addr is %s", config.Accrual)
	InfoLog.Printf("DB addr is %s", config.DatabaseURI)
	handler := NewHandler(config.Accrual, GetCursor(config.DatabaseURI))
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
