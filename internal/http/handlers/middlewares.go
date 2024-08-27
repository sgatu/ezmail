package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	internal_http "github.com/sgatu/ezmail/internal/http"
	"github.com/sgatu/ezmail/internal/http/common"
)

func authorizationMiddleware(appContext *internal_http.AppContext) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("authorization")
			if authHeader != "" {
				if !strings.HasPrefix(authHeader, "Bearer ") {
					unauthorized(w)
					return
				}
				authHeaderSplit := strings.SplitN(authHeader, " ", 2)
				authHeader = authHeaderSplit[1]
				err := appContext.ValidateToken(r.Context(), authHeader)
				if err != nil {
					unauthorized(w)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	common.ErrorResponse(common.BaseError{
		Context:       map[string]string{},
		Message:       "Invalid token",
		ErrIdentifier: "ERR_INVALID_TOKEN",
		Code:          401,
	}, w)
}

func SetupMiddlewares(router *chi.Mux, appContext *internal_http.AppContext) {
	router.Use(authorizationMiddleware(appContext))
}
