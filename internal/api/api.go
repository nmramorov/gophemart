package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nmramorov/gophemart/internal/db"
	"github.com/nmramorov/gophemart/internal/jobmanager"
)

const REQUESTTIMEOUT = 60

type UserRouter struct {
	*chi.Mux
	Cursor *db.Cursor
}

type OrderRouter struct {
	*chi.Mux
	Cursor  *db.Cursor
	Manager *jobmanager.Jobmanager
}

type BalanceRouter struct {
	*chi.Mux
	Cursor *db.Cursor
	// Manager *jobmanager.Jobmanager
}

type Handler struct {
	*chi.Mux
	Cursor *db.Cursor
	// Manager *jobmanager.Jobmanager
}

func NewHandler(cursor *db.Cursor, manager *jobmanager.Jobmanager) *Handler {
	handler := &Handler{
		Mux:    chi.NewMux(),
		Cursor: cursor,
		// Manager: jobmanager.NewJobmanager(cursor, accrualURL),
	}
	handler.Use(GzipHandle)
	handler.Use(handler.CookieHandle)
	handler.Use(middleware.Timeout(REQUESTTIMEOUT * time.Second))

	handler.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handler.RegisterUser)
		r.Post("/login", handler.Login)

		r.Get("/withdrawals", handler.GetWithdrawals)
		r.Get("/balance", handler.GetBalance)
		r.Post("/balance/withdraw", handler.WithdrawMoney)

		OrdersRouter := NewOrdersRouter(cursor, manager)
		r.Mount("/orders", OrdersRouter)
	})

	return handler
}

func NewOrdersRouter(cursor *db.Cursor, manager *jobmanager.Jobmanager) *OrderRouter {
	r := &OrderRouter{
		Mux:     chi.NewMux(),
		Cursor:  cursor,
		Manager: manager,
	}
	r.Post("/", r.UploadOrder)
	r.Get("/", r.GetOrders)
	return r
}
