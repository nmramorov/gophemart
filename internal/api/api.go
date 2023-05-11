package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nmramorov/gophemart/internal/db"
	"github.com/nmramorov/gophemart/internal/jobmanager"
)

const REQUESTTIMEOUT = 60

type Handler struct {
	*chi.Mux
	Cursor  *db.Cursor
	Manager *jobmanager.Jobmanager
}

func NewHandler(accrualURL string, cursor *db.Cursor) *Handler {
	handler := &Handler{
		Mux:     chi.NewMux(),
		Cursor:  cursor,
		Manager: jobmanager.NewJobmanager(cursor, accrualURL),
	}
	handler.Use(GzipHandle)
	handler.Use(handler.CookieHandle)
	handler.Use(middleware.Timeout(REQUESTTIMEOUT * time.Second))

	handler.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handler.RegisterUser)
		r.Post("/login", handler.Login)

		r.Route("/orders", func(r chi.Router) {
			r.Post("/", handler.UploadOrder)
			r.Get("/", handler.GetOrders)
		})

		r.Get("/withdrawals", handler.GetWithdrawals)
		r.Get("/balance", handler.GetBalance)
		r.Post("/balance/withdraw", handler.WithdrawMoney)
		// UserRouter := NewUserRouter(handler)
		// r.Mount("/", UserRouter)

		// OrdersRouter := NewOrdersRouter(handler)
		// r.Mount("/", OrdersRouter)

		// BalanceWithdrawalRouter := NewBalanceWithdrawalsRouter(handler)
		// r.Mount("/", BalanceWithdrawalRouter)
	})

	return handler
}

func NewUserRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()
	r.Post("/register", handler.RegisterUser)
	r.Post("/login", handler.Login)
	return r
}

func NewOrdersRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()
	r.Post("/orders", handler.UploadOrder)
	r.Get("/orders", handler.GetOrders)
	return r
}

func NewBalanceWithdrawalsRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()
	r.Get("/withdrawals", handler.GetWithdrawals)
	r.Get("/balance", handler.GetBalance)
	r.Post("/balance/withdraw", handler.WithdrawMoney)
	return r
}
