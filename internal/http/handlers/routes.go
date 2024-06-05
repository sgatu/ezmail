package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
	"github.com/sgatu/ezmail/internal/http/handlers/auth/login"
	"github.com/sgatu/ezmail/internal/http/handlers/domain"
	"github.com/sgatu/ezmail/internal/http/handlers/register"
)

func secureMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currUser := r.Context().Value(internal_http.CurrentUserKey)
			if currUser == nil {
				common.ErrorResponse(common.UnauthorizedError(), w)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func SetupRoutes(r *chi.Mux, appContext *internal_http.AppContext) {
	login.LoginHandler(appContext, r)
	register.RegisterHandler(appContext, r)
	r.With(secureMiddleware()).Group(func(r chi.Router) {
		domain.DomainHandler(appContext, r)
	})
}
