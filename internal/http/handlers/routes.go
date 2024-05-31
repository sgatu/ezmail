package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
	"github.com/sgatu/ezmail/internal/http/handlers/auth/login"
	"github.com/sgatu/ezmail/internal/http/handlers/register"
)

func secureMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currUser := r.Context().Value(internal_http.CurrentUserKey)
			if currUser == nil {
				common.ReturnError(common.UnauthorizedError(), w)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func SetupRoutes(r *chi.Mux, appContext *internal_http.AppContext) {
	login.LoginHandler(appContext, r)
	register.RegisterHandler(appContext, r)
	r.Group(func(r chi.Router) {
		r.Use(secureMiddleware())
		// after this authorized endpoints
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			currUser := r.Context().Value(internal_http.CurrentUserKey)
			current_user, _ := currUser.(*user.User)
			fmt.Fprintf(w, "Current user - Name: %s, Id: %s", current_user.Name, current_user.Id)
		})
	})
}
