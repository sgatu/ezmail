package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/handlers/auth/login"
	"github.com/sgatu/ezmail/internal/http/handlers/register"
)

func SetupRoutes(r *chi.Mux, appContext *http.AppContext) {
	login.LoginHandler(appContext, r)
	register.RegisterHandler(appContext, r)
}
