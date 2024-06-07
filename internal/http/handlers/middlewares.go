package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sgatu/ezmail/internal/domain/models/user"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
)

func authorizationMiddleware(appContext *internal_http.AppContext) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("authorization")
			if authHeader != "" {
				if !strings.HasPrefix(authHeader, "Bearer ") {
					common.ErrorResponse(common.BaseError{
						Context:       map[string]string{},
						Message:       "Invalid token",
						ErrIdentifier: "ERR_INVALID_TOKEN",
						Code:          400,
					}, w)
					return
				}
				authHeaderSplit := strings.SplitN(authHeader, " ", 2)
				authHeader = authHeaderSplit[1]
				usr := appContext.AuthManager.ValidateToken(r.Context(), authHeader)
				if usr != (*user.User)(nil) {
					r = r.WithContext(context.WithValue(r.Context(), internal_http.CurrentUserKey, usr))
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

func SetupMiddlewares(router *chi.Mux, appContext *internal_http.AppContext) {
	router.Use(authorizationMiddleware(appContext))
}
