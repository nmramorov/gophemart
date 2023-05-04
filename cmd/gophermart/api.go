package main

import (
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*chi.Mux
	Cursor  *Cursor
	Manager *Jobmanager
}

func NewHandler(accrualUrl string, cursor *Cursor) *Handler {
	handler := &Handler{
		Mux:     chi.NewMux(),
		Cursor:  cursor,
		Manager: NewJobmanager(cursor, accrualUrl),
	}
	handler.Post("/api/user/register", handler.RegisterUser)
	handler.Post("/api/user/login", handler.Login)
	handler.Post("/api/user/orders", handler.UploadOrder)
	handler.Get("/api/user/orders", handler.GetOrders)
	handler.Get("/api/user/withdrawals", handler.GetWithdrawals)
	handler.Get("/api/user/balance", handler.GetBalance)
	handler.Post("/api/user/balance/withdraw", handler.WithdrawMoney)

	return handler
}
