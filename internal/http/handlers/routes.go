package handlers

import "github.com/go-chi/chi/v5"

func SetupRoutes(r *chi.Mux) {
	hello := helloHandler{}
	r.Get("/", hello.hello)
}
