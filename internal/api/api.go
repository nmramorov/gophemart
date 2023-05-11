package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nmramorov/gophemart/internal/db"
	"github.com/nmramorov/gophemart/internal/jobmanager"
)

const REQUEST_TIMEOUT = 60

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
	handler.Use(middleware.Timeout(REQUEST_TIMEOUT * time.Second))

	handler.Route("/api/user", func(r chi.Router) {
		userRouter := NewUserRouter(handler)
		r.Mount("/", userRouter)

		ordersRouter := NewOrdersRouter(handler)
		r.Mount("/orders", ordersRouter)

		balanceWithdrawalRouter := NewBalanceWithdrawalsRouter(handler)
		r.Mount("/", balanceWithdrawalRouter)
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
