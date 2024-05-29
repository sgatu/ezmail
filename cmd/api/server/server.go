package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	return router
}
