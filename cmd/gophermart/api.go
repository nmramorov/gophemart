package main

import (
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*chi.Mux
	Cursor *Cursor
}
