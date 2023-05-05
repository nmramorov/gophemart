package main

import (
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*chi.Mux
	Cursor  *Cursor
	DBCursor *Cursor
	Manager *Jobmanager
}

func NewHandler(accrualURL string, cursor *Cursor, dbcursor *Cursor) *Handler {
	handler := &Handler{
		Mux:     chi.NewMux(),
		Cursor:  cursor,
		DBCursor: dbcursor,
		Manager: NewJobmanager(cursor, accrualURL),
	}
	handler.Use(GzipHandle)
	handler.Use(handler.CookieHandle)

	handler.Post("/api/user/register", handler.RegisterUser)
	handler.Post("/api/user/login", handler.Login)
	handler.Post("/api/user/orders", handler.UploadOrder)
	handler.Get("/api/user/orders", handler.GetOrders)
	handler.Get("/api/user/withdrawals", handler.GetWithdrawals)
	handler.Get("/api/user/balance", handler.GetBalance)
	handler.Post("/api/user/balance/withdraw", handler.WithdrawMoney)

	return handler
}
