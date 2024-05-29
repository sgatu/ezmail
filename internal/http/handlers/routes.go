package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/user"
)

func SetupRoutes(r *chi.Mux, userRepo user.UserRepository) {
	hello := helloHandler{userRepo: userRepo}
	r.Get("/", hello.hello)
}
