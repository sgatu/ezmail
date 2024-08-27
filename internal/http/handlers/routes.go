package handlers

import (
	"github.com/go-chi/chi/v5"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/handlers/domain"
	"github.com/sgatu/ezmail/internal/http/handlers/email"
)

/*func secureMiddleware() func(h http.Handler) http.Handler {
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
}*/

func SetupRoutes(r *chi.Mux, appContext *internal_http.AppContext) {
	domain.DomainHandler(appContext, r)
	email.RegisterEmailHandler(appContext, r)
}
