package handlers

import (
	"github.com/go-chi/chi/v5"
	internalhttp "github.com/sgatu/ezmail/internal/http"
)

func SetupRoutes(r *chi.Mux, appContext *internalhttp.AppContext) {
	RegisterDomainHandler(appContext, r)
	RegisterEmailHandler(appContext, r)
	RegisterTemplateHandler(appContext, r)
}
